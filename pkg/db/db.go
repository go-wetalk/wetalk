package db

import (
	"appsrv/pkg/config"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	sql, _ := q.FormattedQuery()
	fmt.Println(sql)
	return nil
}

var db *pg.DB

// ProvideSingleton provides singleton DB instance.
func ProvideSingleton() *pg.DB {
	if db == nil {
		c := config.ProvideSingleton()

		db = pg.Connect(&pg.Options{
			Addr:     c.DB.Addr,
			User:     c.DB.User,
			Password: c.DB.Secret,
			Database: c.DB.Name,
		})

		if c.Debug {
			db.AddQueryHook(dbLogger{})
		}

		if _, err := db.Exec("select 1"); err != nil {
			panic(err)
		}
	}

	return db
}
