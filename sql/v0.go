package sql

import (
	"time"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func init() {
	Setup(0, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			&version{},
		)

		return nil
	})
}

type version struct {
	ID      uint
	Code    int       `pg:",unique,notnull"`
	Created time.Time `pg:",default:now()"`
}
