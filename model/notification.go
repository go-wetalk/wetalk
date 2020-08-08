package model

import "appsrv/pkg/db"

// Notification 系统通知
type Notification struct {
	ID      uint
	RecvID  uint
	Content string
	HasRead bool `pg:",default:false"`

	db.TimeUpdate
}
