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

		db.Insert(
			&siteConfig{
				Key:     "domain",
				Value:   "devto.icu",
				Comment: "域名，用于生成对外链接，如：devto.icu",
			},
			&siteConfig{
				Key:     "name",
				Value:   "DevToICU",
				Comment: "社区名称",
			},
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
