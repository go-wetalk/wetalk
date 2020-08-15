package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/db"
	"appsrv/pkg/out"
	"appsrv/schema"
	"appsrv/service"
	"net/http"

	"github.com/kataras/muxie"
	"github.com/spf13/cast"
)

var Comment = new(comment)

type comment struct{}

// CreateTopicComment 发表帖子评论
func (comment) CreateTopicComment(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err(401, "您还没有登录，请先登录"))
		return
	}

	var input schema.TopicCommentCreation
	err = muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	if err = input.Validate(); err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	c, err := service.Comment.CreateTopicComment(db.DB, u, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(c))
}

func (comment) ListByFilter(w http.ResponseWriter, r *http.Request) {
	input := schema.CommentFilter{}
	input.TopicID = cast.ToUint(r.URL.Query().Get("tid"))
	input.Page = cast.ToInt(r.URL.Query().Get("p"))
	if input.Page < 1 {
		input.Page = 1
	}
	input.Size = cast.ToInt(r.URL.Query().Get("s"))
	if input.Size < 1 {
		input.Size = 20
	}

	cs, err := service.Comment.FindByFilterInput(db.DB, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(cs))
}
