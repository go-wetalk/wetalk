package model

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/db"

	"go.uber.org/zap"
)

type Admin struct {
	ID       uint
	Name     string
	Password string   `json:"-"`
	RoleKeys []string `pg:",array"`

	db.TimeUpdate
}

func (a *Admin) RoleList() (roles []Role) {
	q := "select a.* from roles a join admin_roles b on b.admin_id = ? and b.role_id = a.id"
	_, err := db.DB.Query(&roles, q, a.ID)
	if err != nil {
		bog.Error("Admin.RoleList", zap.Error(err))
	}
	return
}
