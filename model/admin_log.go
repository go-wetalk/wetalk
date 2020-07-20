package model

import "time"

type AdminLog struct {
	ID        uint
	AdminID   uint
	AdminName string
	Event     string
	Intro     string
	IP        string
	UA        string
	Ref       string
	Created   time.Time `pg:",default:now()"`
}
