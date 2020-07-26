package app

import (
	"appsrv/model"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"net/http"

	"github.com/kataras/muxie"
	"go.uber.org/zap"
)

type Role struct{}

func (Role) List(w http.ResponseWriter, r *http.Request) {
	var rs = []model.Role{}
	err := db.DB.Model(&rs).Select()
	if err != nil {
		bog.Error("Role.List", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, &rs)
}

func (Role) Create(w http.ResponseWriter, r *http.Request) {

}
