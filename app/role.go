package app

import (
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"net/http"

	"github.com/kataras/muxie"
	"go.uber.org/zap"
)

type Role struct {
	ID     uint
	Key    string
	Name   string
	Intro  string
	Admins []Admin `pg:"many2many:admin_roles,joinFK:role_id" json:"-"`
}

func (Role) CheckRole(names ...string) muxie.Wrapper {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var u Admin
			err := auth.GetUser(r, &u)
			if err != nil {
				w.WriteHeader(401)
				return
			}

			roles := u.RoleList()
			yes := false
		TOP:
			for _, role := range roles {
				for _, name := range names {
					if role.Key == name {
						yes = true
						break TOP
					}
				}
			}

			if yes {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(403)
			}
		})
	}
}

func (Role) List(w http.ResponseWriter, r *http.Request) {
	var rs = []Role{}
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