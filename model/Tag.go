package model

import "appsrv/pkg/db"

type Tag struct {
	ID        uint
	Name      string
	CreatorID uint
	Creator   User `pg:",foreginKey:creator_id"`

	db.TimeUpdate
}
