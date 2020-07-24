package model

import "appsrv/pkg/db"

type Comment struct {
	ID        uint
	TopicID   uint
	UserID    uint
	CommentID uint `pg:",default:0"`
	Content   string

	db.TimeUpdate

	User    *User
	Topic   *Topic
	Comment *Comment
}
