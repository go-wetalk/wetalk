package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/db"
	"appsrv/pkg/errors"
	"appsrv/schema"
	"appsrv/service"
	"net/http"
	"strings"

	"github.com/kataras/muxie"
	"github.com/spf13/cast"
)

// Topic 话题
type Topic struct{}

// List 取出话题列表
func (Topic) List(w http.ResponseWriter, r *http.Request) {
	input := schema.TopicListInput{}
	input.Size = 20
	if p := r.URL.Query().Get("p"); p != "" {
		input.Page = cast.ToUint(p)
	}
	if t := r.URL.Query().Get("t"); t != "" {
		input.Tag = strings.TrimSpace(t)
	}

	ts, _ := service.Topic.ListWithRankByScore(db.DB, input)
	muxie.Dispatch(w, muxie.JSON, ts)
}

// Create 创建话题
func (Topic) Create(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		w.WriteHeader(401)
		muxie.Dispatch(w, muxie.JSON, errors.New(401, "请登录"))
		return
	}

	var input schema.TopicCreateInput
	err = muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		w.WriteHeader(429)
		muxie.Dispatch(w, muxie.JSON, errors.ErrBodyBind)
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	if err = input.Validate(); err != nil {
		w.WriteHeader(400)
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	t, err := service.Topic.Create(db.DB, u, input)
	if err != nil {
		w.WriteHeader(500)
		muxie.Dispatch(w, muxie.JSON, errors.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, t)
}

// Find 查看话题详情
func (Topic) Find(w http.ResponseWriter, r *http.Request) {
	topicID := cast.ToUint(muxie.GetParam(w, "topicID"))
	t, err := service.Topic.FindByID(db.DB, topicID)
	if err != nil {
		w.WriteHeader(500)
		muxie.Dispatch(w, muxie.JSON, errors.Err500)
	}

	muxie.Dispatch(w, muxie.JSON, t)
}
