package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"net/http"
	"time"

	"github.com/kataras/muxie"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type Task struct{}

func (Task) List(w http.ResponseWriter, r *http.Request) {
	var ts = []model.Task{}
	_ = db.DB.Model(&ts).Order("id ASC").Select()
	muxie.Dispatch(w, muxie.JSON, &ts)
}

func (Task) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		model.Task
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Task.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	err = db.DB.Insert(&in.Task)
	if err != nil {
		bog.Error("Task.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (Task) Update(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name    string
		Content string
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Task.Update", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	var t model.Task
	err = db.DB.Model(&t).Where("id = ?", muxie.GetParam(w, "TaskID")).First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	_, err = db.DB.Model(&t).WherePK().Set("name = ?", in.Name).Set("content = ?", in.Content).Update()
	if err != nil {
		bog.Error("Task.Update", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (Task) Delete(w http.ResponseWriter, r *http.Request) {
	m := model.Task{}
	err := db.DB.Model(&m).Where("id = ?", muxie.GetParam(w, "taskID")).First()
	if err == nil {
		_, err = db.DB.Model(&m).WherePK().Delete()
		if err != nil {
			bog.Error("Task.Delete", zap.Error(err))
			w.WriteHeader(500)
			return
		}
	}

	w.WriteHeader(204)
}

func (Task) AppList(w http.ResponseWriter, r *http.Request) {
	ts := []model.Task{}
	now := time.Now()
	count, err := db.DB.Model(&ts).
		Where("(\"begin\" IS NULL OR \"begin\" <= ?) AND (\"end\" IS NULL OR \"end\" >= ?)", now, now).
		Order("seq DESC").
		SelectAndCount()
	if err != nil {
		bog.Error("Task.AppList", zap.Error(err))
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
		if !t.AvailableTo(&u) {
			continue
		}

		i := outItem{
			ID:         t.ID,
			Name:       t.Name,
			Intro:      t.Intro,
			Button:     t.Fulfilled(&u) && !t.Got(&u),
			StatusText: t.StatusText(&u),
		}
		out.Data = append(out.Data, i)
	}
	muxie.Dispatch(w, muxie.JSON, &out)
}

func (Task) AppTaskLogCreate(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var t model.Task
	kword := muxie.GetParam(w, "taskID")
	err = db.DB.Model(&t).Where("id = ? or factor = ?", cast.ToInt(kword), kword).Order("id DESC").First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if !t.AvailableTo(&u) {
		w.WriteHeader(404)
		return
	}

	if !t.Fulfilled(&u) {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	if t.Got(&u) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	_, err = t.Confirm(&u)
	if err != nil {
		bog.Error("Task.AppTaskLogCreate", zap.Error(err), zap.Uint("TaskID", t.ID), zap.Uint("UserID", u.ID))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}
