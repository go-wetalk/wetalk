package sql

import (
	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func init() {
	Setup(6, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			&siteConfig{},
		)
		return nil
	})
}

type siteConfig struct {
	Key     string `pg:",unique,notnull"`
	Value   string
	Version uint `pg:",default:0"`
	Comment string
}
