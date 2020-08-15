package admin

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"appsrv/pkg/out"
	"appsrv/service"
	"net/http"

	"github.com/kataras/muxie"
	"go.uber.org/zap"
)

type User struct{}

func (User) Create(w http.ResponseWriter, r *http.Request) {
	input := service.UserCreateInput{}
	err := muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	u, err := service.User{}.Create(db.DB, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, u)
}

func (User) List(w http.ResponseWriter, r *http.Request) {
	users, err := service.User{}.List(db.DB)
	if err != nil {
		bog.Error("User.List", zap.Error(err))
	}
	muxie.Dispatch(w, muxie.JSON, users)
}
