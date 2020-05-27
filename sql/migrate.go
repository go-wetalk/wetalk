package sql

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"sort"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"go.uber.org/zap"
)

var ms = map[int]func(db *pg.DB, l *zap.Logger) error{}

func Setup(v int, m func(db *pg.DB, l *zap.Logger) error) {
	ms[v] = m
}

func Run() {
	ints := make([]int, 0)
	for i := range ms {
		ints = append(ints, i)
	}

	sort.Ints(ints)

	for _, v := range ints {
		m := ms[v]
		if checkV(v) {
			bog.Info("migating", zap.Int("v", v))
			err := m(db.DB, bog.Log)
			if err != nil {
				bog.Fatal("migrating fails", zap.Error(err))
			}

			updateV(v)
			bog.Info("migrated", zap.Int("v", v))
		} else {
			bog.Info("migration skipped", zap.Int("v", v))
		}
	}
}

func checkErr(err error) {
	if err != nil {
		bog.Error("migration error", zap.Error(err))
	}
}

func createTable(tables ...interface{}) {
	for _, ptr := range tables {
		err := db.DB.CreateTable(ptr, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		checkErr(err)
	}
}

func checkV(v int) bool {
	if v == 0 {
		return true
	}

	var max int
	db.DB.Model((*version)(nil)).ColumnExpr("MAX(version.code)").Select(&max)
	return v > max
}

func updateV(v int) error {
	return db.DB.Insert(&version{
		Code: v,
	})
}
