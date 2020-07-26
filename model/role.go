package model

import "appsrv/pkg/db"

type Role struct {
	ID    uint
	Key   string `pg:",unique"`
	Name  string
	Intro string

	db.TimeUpdate
}
