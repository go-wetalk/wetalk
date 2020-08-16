package model

import (
	"appsrv/pkg/db"
	"errors"

	"github.com/go-pg/pg/v9"
)

// TaskLog 任务日志
type TaskLog struct {
	ID        uint
	TaskID    uint
	UserID    uint
	Bonus     int    // 奖励类型
	BonusNum  int    `pg:",default:0"` // 奖励数额
	Factor    string // 任务因素名称
	FactorNum int    `pg:",default:0"` // 完成需要的因素数量

	db.TimeUpdate
}

// var _ pg.AfterInsertHook = (*TaskLog)(nil)

// func (l *TaskLog) AfterInsert(ctx context.Context) error {
// 	return l.DispatchBonus()
// }

func (l *TaskLog) DispatchBonus(db *pg.DB) error {
	switch l.Bonus {
	case BonusCoin:
		var bal int
		db.Model(&User{ID: l.UserID}).WherePK().Column("coin").Select(pg.Scan(&bal))
		cl := CoinLog{
			UserID:   l.UserID,
			Source:   "tasks",
			SourceID: l.ID,
			Value:    l.BonusNum,
			Balance:  bal + l.BonusNum,
		}
		return db.Insert(&cl)
	}
	return errors.New("undefined task bonus")
}
