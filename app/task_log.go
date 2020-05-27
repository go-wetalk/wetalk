package app

import (
	"appsrv/pkg/db"
	"context"
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

var _ pg.AfterInsertHook = (*TaskLog)(nil)

func (l *TaskLog) AfterInsert(ctx context.Context) error {
	return l.DispatchBonus()
}

func (l *TaskLog) DispatchBonus() error {
	switch l.Bonus {
	case BonusCoin:
		var bal int
		db.DB.Model(&User{ID: l.UserID}).WherePK().Column("coin").Select(pg.Scan(&bal))
		cl := CoinLog{
			UserID:   l.UserID,
			Source:   "tasks",
			SourceID: l.ID,
			Value:    l.BonusNum,
			Balance:  bal + l.BonusNum,
		}
		return db.DB.Insert(&cl)
	}
	return errors.New("undefined task bonus")
}
