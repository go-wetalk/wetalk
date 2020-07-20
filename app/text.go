package app

import (
	"appsrv/model"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"net/http"

	"github.com/kataras/muxie"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type Text struct{}

func (Text) List(w http.ResponseWriter, r *http.Request) {
	var ts = []model.Text{}
	_ = db.DB.Model(&ts).Order("id ASC").Select()
	muxie.Dispatch(w, muxie.JSON, &ts)
}

func (Text) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		model.Text
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Text.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	err = db.DB.Insert(&in.Text)
	if err != nil {
		bog.Error("Text.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (Text) Update(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name    string
		Content string
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Text.Update", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	var t model.Text
	err = db.DB.Model(&t).Where("id = ?", muxie.GetParam(w, "textID")).First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	_, err = db.DB.Model(&t).WherePK().Set("name = ?", in.Name).Set("content = ?", in.Content).Update()
	if err != nil {
		bog.Error("Text.Update", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (Text) AppView(w http.ResponseWriter, r *http.Request) {
	var t model.Text
	textID := muxie.GetParam(w, "textID")
	err := db.DB.Model(&t).Where("id = ? OR slot_name = ?", cast.ToUint(textID), textID).First()
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
