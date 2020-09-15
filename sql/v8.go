package sql

import (
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

func init() {
	Setup(8, func(db *pg.DB, l *zap.Logger) error {
		db.Model(
			&rule{
				Path:       "/users/*",
				Method:     "{GET}",
				Authorized: []string{"*"},
			},
		).Insert()
		db.Model(
			&rule{
				Path:       "/notifications",
				Method:     "{GET}",
				Authorized: []string{"*"},
			},
		).Insert()
		db.Model(
			&rule{
				Path:       "/notifications/*",
				Method:     "{DELETE}",
				Authorized: []string{"*"},
			},
		).Insert()
		return nil
	})
}
