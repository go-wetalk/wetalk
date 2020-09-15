package sql

import (
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

func init() {
	Setup(9, func(db *pg.DB, l *zap.Logger) error {
		db.Model(
			&rule{
				Path:       "/profile",
				Method:     "{GET}",
				Authorized: []string{"*"},
			},
		).Insert()
		db.Model(
			&rule{
				Path:       "/profile/*",
				Method:     "{PUT}",
				Authorized: []string{"*"},
			},
		).Insert()
		return nil
	})
}
