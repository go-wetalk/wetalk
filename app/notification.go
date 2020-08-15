package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"appsrv/pkg/out"
	"appsrv/schema"
	"appsrv/service"
	"net/http"

	"github.com/kataras/muxie"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// Notification 消息
var Notification = &notification{}

type notification struct{}

func (notification) List(w http.ResponseWriter, r *http.Request) {
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

	ret, err := service.Notification.FindForUser(db.DB, &u, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(ret))
}

func (notification) MarkRead(w http.ResponseWriter, r *http.Request) {
	notifyID := cast.ToUint(muxie.GetParam(w, "notificationID"))

	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	err = service.Notification.MarkAsRead(db.DB, &u, notifyID)
	if err != nil {
		bog.Error("notification.MarkRead", zap.Error(err))
	}

	muxie.Dispatch(w, muxie.JSON, out.Err(204, "操作成功"))
}
