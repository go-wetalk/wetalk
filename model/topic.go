package model

import "appsrv/pkg/db"

type Topic struct {
	ID      uint
	UserID  uint
	Title   string
	Content string
	Tags    []Tag `pg:",many2many:topic_tag"`

	db.TimeUpdate

	User *User
}
