package sql

import (
	"sort"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/prometheus/common/log"
	"go.uber.org/zap"
)

var ms = map[int]func(db *pg.DB, l *zap.Logger) error{}

func Setup(v int, m func(db *pg.DB, l *zap.Logger) error) {
	ms[v] = m
}

func Run(db *pg.DB, log *zap.Logger) {
	ints := make([]int, 0)
	for i := range ms {
		ints = append(ints, i)
	}

	sort.Ints(ints)

	for _, v := range ints {
		m := ms[v]
		if checkV(db, v) {
			log.Info("migating", zap.Int("v", v))
			err := m(db, log)
			if err != nil {
				log.Fatal("migrating fails", zap.Error(err))
			}

			updateV(db, v)
			log.Info("migrated", zap.Int("v", v))
		} else {
			log.Info("migration skipped", zap.Int("v", v))
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Error("migration error", zap.Error(err))
	}
}

func createTable(db *pg.DB, tables ...interface{}) {
	for _, ptr := range tables {
		err := db.CreateTable(ptr, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		checkErr(err)
	}
}

func checkV(db *pg.DB, v int) bool {
	if v == 0 {
		return true
	}

	var max int
	db.Model((*version)(nil)).ColumnExpr("MAX(version.code)").Select(&max)
	return v > max
}

func updateV(db *pg.DB, v int) error {
	return db.Insert(&version{
		Code: v,
	})
}
