package sql

import (
	"appsrv/pkg/db"
	"time"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

func init() {
	Setup(4, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			db,
			&task{},
			&taskLog{},
			&coinLog{},
		)

		return nil
	})
}

type task struct {
	ID        uint
	Name      string
	Intro     string
	Bonus     int    // 奖励类型
	BonusNum  int    `pg:",default:0"` // 奖励数额
	Factor    string // 任务因素名称
	FactorNum int    `pg:",default:0"` // 完成需要的因素数量

	Daily   bool       `pg:",default:false"` // 标记是否是每日任务
	Cooling uint       `pg:",default:0"`     // 标记任务刷新间隔，0表示无刷新间隔
	Times   uint       `pg:",default:0"`     // 标记任务限制次数，0表示无限制
	Begin   *time.Time // 起始时间
	End     *time.Time // 截止时间

	Seq int `pg:",default:0"`

	db.TimeUpdate
}

type taskLog struct {
	ID        uint
	TaskID    uint
	UserID    uint
	Bonus     int    // 奖励类型
	BonusNum  int    `pg:",default:0"` // 奖励数额
	Factor    string // 任务因素名称
	FactorNum int    `pg:",default:0"` // 完成需要的因素数量

	db.TimeUpdate
}

type coinLog struct {
	ID       uint
	UserID   uint
	Source   string // 变动因素
	SourceID uint
	Value    int
	Balance  int

	db.TimeUpdate
}
