package sql

import (
	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func init() {
	Setup(8, func(db *pg.DB, l *zap.Logger) error {
		db.Insert(
			&rule{
				Path:       "/users/*",
				Method:     "{GET}",
				Authorized: []string{"*"},
			},
		)
		db.Insert(
			&rule{
				Path:       "/notifications",
				Method:     "{GET}",
				Authorized: []string{"*"},
			},
		)
		db.Insert(
			&rule{
				Path:       "/notifications/*",
				Method:     "{DELETE}",
				Authorized: []string{"*"},
			},
		)
		return nil
	})
}
