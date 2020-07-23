package db

import (
	"appsrv/pkg/config"
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
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

var DB *pg.DB

func InitConn(c config.DBConfig, debug bool) error {
	DB = pg.Connect(&pg.Options{
		Addr:     c.Addr,
		User:     c.User,
		Password: c.Secret,
		Database: c.Name,
	})

	if debug {
		DB.AddQueryHook(dbLogger{})
	}

	_, err := DB.Exec("select 1")
	return err
}
