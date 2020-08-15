package sql

import (
	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func init() {
	Setup(9, func(db *pg.DB, l *zap.Logger) error {
		db.Insert(
			&rule{
				Path:       "/profile",
				Method:     "{GET}",
				Authorized: []string{"*"},
			},
		)
		db.Insert(
			&rule{
				Path:       "/profile/*",
				Method:     "{PUT}",
				Authorized: []string{"*"},
			},
		)
		return nil
	})
}
