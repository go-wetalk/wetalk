package sql

import (
	"appsrv/pkg/db"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

func init() {
	Setup(5, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			db,
			&topic{},
			&comment{},
		)

		db.Model([]*rule{
			&rule{
				Path:       "/topics",
				Method:     "{POST}",
				Authorized: []string{"*"},
			},
			&rule{
				Path:       "/topics/*",
				Method:     "{PUT,DELETE}",
				Authorized: []string{"*"},
			},
			&rule{
				Path:        "/topics/*",
				Method:      "{GET}",
				AllowAnyone: true,
			},
			&rule{
				Path:       "/comments",
				Method:     "{POST}",
				Authorized: []string{"*"},
			},
			&rule{
				Path:        "/comments",
				Method:      "{GET}",
				AllowAnyone: true,
			},
		}).Insert()
		return nil
	})
}

type topic struct {
	ID      uint
	UserID  uint
	Title   string
	Content string
	Tags    []string `pg:",array"`

	db.TimeUpdate
}

type comment struct {
	ID      uint
	TopicID uint
	UserID  uint
	Content string

	db.TimeUpdate
}
