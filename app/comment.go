//+build wireinject

package app

import (
	"appsrv/model"
	"appsrv/pkg"
	"appsrv/pkg/auth"
	"appsrv/pkg/config"
	"appsrv/pkg/out"
	"appsrv/pkg/runtime"
	"appsrv/schema"
	"appsrv/service"
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/google/wire"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type Comment struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}

func (v *Comment) RegisterRoute(m muxie.SubMux) {
	m.Handle("/comments", muxie.Methods().
		HandleFunc(http.MethodGet, v.ListByFilter).
		HandleFunc(http.MethodPost, v.CreateTopicComment))
}

func NewCommentController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Comment), "*"),
		wire.Bind(new(runtime.Controller), new(*Comment)),
	)
	return nil
}

// CreateTopicComment 发表帖子评论
func (v Comment) CreateTopicComment(w http.ResponseWriter, r *http.Request) {
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

	c, err := service.Comment.CreateTopicComment(v.db, u, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(c))
}

func (v Comment) ListByFilter(w http.ResponseWriter, r *http.Request) {
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

	cs, err := service.Comment.FindByFilterInput(v.db, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(cs))
}
