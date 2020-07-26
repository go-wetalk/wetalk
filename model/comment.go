package model

import "appsrv/pkg/db"

type Comment struct {
	ID      uint
	TopicID uint
	UserID  uint
	Content string

	db.TimeUpdate

	User  *User
	Topic *Topic
}
