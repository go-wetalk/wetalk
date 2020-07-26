package model

import (
	"appsrv/pkg/db"
)

type Admin struct {
	ID       uint
	Name     string
	Password string   `json:"-"`
	RoleKeys []string `pg:",array"`

	db.TimeUpdate
}
