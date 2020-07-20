package sql

import (
	"appsrv/pkg/db"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func init() {
	Setup(5, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			&topic{},
		)
		return nil
	})
}

type topic struct {
	ID      uint
	UserID  uint
	Title   string
	Content string
	Tags    []tag `pg:",many2many:topic_tag"`

	db.TimeUpdate
}

type tag struct {
	ID        uint
	Name      string
	CreatorID uint

	db.TimeUpdate
}
