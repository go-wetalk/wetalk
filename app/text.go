package app

import (
	"appsrv/model"
	"appsrv/pkg/config"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type Text struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}

func (v *Text) RegisterRoute(m muxie.SubMux) {
	m.Handle("/texts/:textID", muxie.Methods().HandleFunc(http.MethodGet, v.AppView))
}

func (v Text) List(w http.ResponseWriter, r *http.Request) {
	var ts = []model.Text{}
	_ = v.db.Model(&ts).Order("id ASC").Select()
	muxie.Dispatch(w, muxie.JSON, &ts)
}

func (v Text) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		model.Text
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		v.log.Error("Text.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	_, err = v.db.Model(&in.Text).Insert()
	if err != nil {
		v.log.Error("Text.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (v Text) Update(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name    string
		Content string
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		v.log.Error("Text.Update", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	var t model.Text
	err = v.db.Model(&t).Where("id = ?", muxie.GetParam(w, "textID")).First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	_, err = v.db.Model(&t).WherePK().Set("name = ?", in.Name).Set("content = ?", in.Content).Update()
	if err != nil {
		v.log.Error("Text.Update", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (v Text) AppView(w http.ResponseWriter, r *http.Request) {
	var t model.Text
	textID := muxie.GetParam(w, "textID")
	err := v.db.Model(&t).Where("id = ? OR slot_name = ?", cast.ToUint(textID), textID).First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	var out struct {
		Name    string
		Content string
	}
	out.Name = t.Name
	out.Content = t.Content
	muxie.Dispatch(w, muxie.JSON, &out)
}
