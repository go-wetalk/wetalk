package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/config"
	"net/http"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type Task struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}

func (v *Task) RegisterRoute(m muxie.SubMux) {
	m.Handle("/tasks", muxie.Methods().
		HandleFunc(http.MethodGet, v.AppList))
	m.Handle("/tasks/:taskID/bonus", muxie.Methods().
		HandleFunc(http.MethodPost, v.AppTaskLogCreate))
}

func (v Task) List(w http.ResponseWriter, r *http.Request) {
	var ts = []model.Task{}
	_ = v.db.Model(&ts).Order("id ASC").Select()
	muxie.Dispatch(w, muxie.JSON, &ts)
}

func (v Task) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		model.Task
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		v.log.Error("Task.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	err = v.db.Insert(&in.Task)
	if err != nil {
		v.log.Error("Task.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (v Task) Update(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name    string
		Content string
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		v.log.Error("Task.Update", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	var t model.Task
	err = v.db.Model(&t).Where("id = ?", muxie.GetParam(w, "TaskID")).First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	_, err = v.db.Model(&t).WherePK().Set("name = ?", in.Name).Set("content = ?", in.Content).Update()
	if err != nil {
		v.log.Error("Task.Update", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (v Task) Delete(w http.ResponseWriter, r *http.Request) {
	m := model.Task{}
	err := v.db.Model(&m).Where("id = ?", muxie.GetParam(w, "taskID")).First()
	if err == nil {
		_, err = v.db.Model(&m).WherePK().Delete()
		if err != nil {
			v.log.Error("Task.Delete", zap.Error(err))
			w.WriteHeader(500)
			return
		}
	}

	w.WriteHeader(204)
}

func (v Task) AppList(w http.ResponseWriter, r *http.Request) {
	ts := []model.Task{}
	now := time.Now()
	count, err := v.db.Model(&ts).
		Where("(\"begin\" IS NULL OR \"begin\" <= ?) AND (\"end\" IS NULL OR \"end\" >= ?)", now, now).
		Order("seq DESC").
		SelectAndCount()
	if err != nil {
		v.log.Error("Task.AppList", zap.Error(err))
	}

	var u model.User
	auth.GetUser(r, &u)

	type outItem struct {
		ID         uint
		Name       string
		Intro      string
		Button     bool // 按钮形式就表示已完成未领取
		StatusText string
	}

	var out = struct {
		Count int
		Page  int
		Data  []outItem
	}{
		Count: count,
		Page:  1,
		Data:  []outItem{},
	}
	for _, t := range ts {
		if t.AvailableTo(v.db, &u) != nil {
			continue
		}

		i := outItem{
			ID:         t.ID,
			Name:       t.Name,
			Intro:      t.Intro,
			Button:     t.Fulfilled(&u) && !t.Got(v.db, &u),
			StatusText: t.StatusText(v.db, &u),
		}
		out.Data = append(out.Data, i)
	}
	muxie.Dispatch(w, muxie.JSON, &out)
}

func (v Task) AppTaskLogCreate(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var t model.Task
	kword := muxie.GetParam(w, "taskID")
	err = v.db.Model(&t).Where("id = ? or factor = ?", cast.ToInt(kword), kword).Order("id DESC").First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if t.AvailableTo(v.db, &u) != nil {
		w.WriteHeader(404)
		return
	}

	if !t.Fulfilled(&u) {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	if t.Got(v.db, &u) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	_, err = t.Confirm(v.db, &u)
	if err != nil {
		v.log.Error("Task.AppTaskLogCreate", zap.Error(err), zap.Uint("TaskID", t.ID), zap.Uint("UserID", u.ID))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}
