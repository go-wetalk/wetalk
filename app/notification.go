package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/config"
	"appsrv/pkg/out"
	"appsrv/schema"
	"appsrv/service"
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type Notification struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}

func (v *Notification) RegisterRoute(m muxie.SubMux) {
	m.Handle("/notifications", muxie.Methods().HandleFunc(http.MethodGet, v.List))
	m.Handle("/notifications/:notificationID", muxie.Methods().HandleFunc(http.MethodDelete, v.MarkRead))
}

func (v Notification) List(w http.ResponseWriter, r *http.Request) {
	input := schema.Paginate{}
	input.Size = cast.ToInt(r.URL.Query().Get("s"))
	if input.Size < 1 {
		input.Size = 20
	}
	input.Page = cast.ToInt(r.URL.Query().Get("p"))
	if input.Page < 1 {
		input.Page = 1
	}

	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	ret, err := service.Notification.FindForUser(v.db, &u, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(ret))
}

func (v Notification) MarkRead(w http.ResponseWriter, r *http.Request) {
	notifyID := cast.ToUint(muxie.GetParam(w, "notificationID"))

	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	err = service.Notification.MarkAsRead(v.db, &u, notifyID)
	if err != nil {
		v.log.Error("notification.MarkRead", zap.Error(err))
	}

	muxie.Dispatch(w, muxie.JSON, out.Err(204, "操作成功"))
}
