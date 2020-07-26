package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/db"
	"appsrv/pkg/errors"
	"appsrv/schema"
	"appsrv/service"
	"net/http"

	"github.com/kataras/muxie"
)

type Comment struct{}

// CreateTopicComment 发表帖子评论
func (Comment) CreateTopicComment(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		w.WriteHeader(401)
		muxie.Dispatch(w, muxie.JSON, errors.New(401, "您还没有登录，请先登录"))
		return
	}

	var input schema.TopicCommentCreation
	err = muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		w.WriteHeader(429)
		muxie.Dispatch(w, muxie.JSON, errors.ErrBodyBind)
		return
	}

	if err = input.Validate(); err != nil {
		w.WriteHeader(400)
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	c, err := service.Comment{}.CreateTopicComment(db.DB, u, input)
	if err != nil {
		w.WriteHeader(err.(errors.JSONError).Code)
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, c)
}
