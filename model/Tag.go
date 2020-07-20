package model

import "appsrv/pkg/db"

type Tag struct {
	ID        uint
	Name      string
	CreatorID uint

	db.TimeUpdate
}
