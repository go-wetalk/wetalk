package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/config"
	"appsrv/pkg/out"
	"appsrv/schema"
	"appsrv/service"
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// Topic 话题
type Topic struct {
	db           *pg.DB
	log          *zap.Logger
	mc           *minio.Client
	conf         *config.ServerConfig
	topicService *service.Topic
}

func (v *Topic) RegisterRoute(m muxie.SubMux) {
	m.Handle("/topics", muxie.Methods().
		HandleFunc(http.MethodGet, v.List).
		HandleFunc(http.MethodPost, v.Create))

	m.Handle("/topics/:topicID", muxie.Methods().
		HandleFunc(http.MethodGet, v.Find))
}

// List 取出话题列表
func (v *Topic) List(w http.ResponseWriter, r *http.Request) {
	input := schema.TopicListInput{}
	input.Size = cast.ToInt(r.URL.Query().Get("s"))
	if input.Size < 1 {
		input.Size = 20
	}
	input.Page = cast.ToInt(r.URL.Query().Get("p"))
	if input.Page < 1 {
		input.Page = 1
	}
	if t := r.URL.Query().Get("t"); t != "" {
		input.Tag = strings.TrimSpace(t)
	}

	ts, _ := v.topicService.ListWithRankByScore(input)
	muxie.Dispatch(w, muxie.JSON, out.Data(ts))
}

// Create 创建话题
func (v *Topic) Create(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err401)
		return
	}

	var input schema.TopicCreateInput
	err = muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	if err = input.Validate(); err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err(400, err.Error()))
		return
	}

	t, err := v.topicService.Create(u, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(t))
}

// Find 查看话题详情
func (v *Topic) Find(w http.ResponseWriter, r *http.Request) {
	topicID := cast.ToUint(muxie.GetParam(w, "topicID"))
	t, err := v.topicService.FindByID(topicID)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(t))
}
