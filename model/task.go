package model

import (
	"appsrv/pkg/db"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/jinzhu/now"
)

const (
	// BonusCoin 积分奖励
	BonusCoin = 1
)

const (
	// FactorCheckIn 签到
	FactorCheckIn = "checkin"
	// FactorUserCreate 邀请注册
	FactorUserCreate = "user:create"
)

// Task 任务表
type Task struct {
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

// AvailableTo 判断该任务对于给定用户是否可完成
func (t *Task) AvailableTo(db *pg.DB, u *User) error {
	n := time.Now()
	if t.Begin != nil && t.Begin.After(n) {
		return errors.New("任务未开始")
	}
	if t.End != nil && t.End.Before(n) {
		return errors.New("任务已结束")
	}

	l := TaskLog{}
	err := db.Model(&l).Where("task_id = ? AND user_id = ?", t.ID, u.ID).Order("id DESC").First()
	if err != nil {
		return nil
	}

	// 每日刷新的任务，按天计算，不按小时计算
	if t.Daily && l.Created.Before(now.BeginningOfDay()) {
		return nil
	}

	if t.Cooling > 0 && l.Created.Add(time.Duration(t.Cooling)*time.Second).Before(n) {
		return nil
	}

	count, _ := db.Model((*TaskLog)(nil)).Where("task_id = ? AND user_id = ?", t.ID, u.ID).Count()
	if t.Times > 0 && int(t.Times) > count {
		return nil
	}

	return errors.New("条件不满足")
}

func (t *Task) StatusText(db *pg.DB, u *User) string {
	switch t.Bonus {
	case BonusCoin:
		if t.Fulfilled(u) {
			if t.Got(db, u) {
				return "已领取"
			} else {
				return fmt.Sprintf("积分 +%d", t.BonusNum)
			}
		} else {
			return fmt.Sprintf("%d/%d", t.Step(u), t.FactorNum)
		}
	}

	return "迷"
}

func (t *Task) Step(u *User) int {
	switch t.Factor {
	case FactorCheckIn:
		return 1
	}
	return 0
}

func (t *Task) Fulfilled(u *User) bool {
	return t.Step(u) == t.FactorNum
}

func (t *Task) Got(db *pg.DB, u *User) bool {
	count, _ := db.Model((*TaskLog)(nil)).
		Where("task_id = ? AND user_id = ?", t.ID, u.ID).
		Where("created BETWEEN ? AND ?", now.BeginningOfDay(), now.EndOfDay()).
		Count()
	return count > 0
}

func (t *Task) Confirm(db *pg.DB, u *User) (*TaskLog, error) {
	l := TaskLog{
		TaskID:    t.ID,
		UserID:    u.ID,
		Bonus:     t.Bonus,
		BonusNum:  t.BonusNum,
		Factor:    t.Factor,
		FactorNum: t.FactorNum,
	}
	_, err := db.Model(&l).Insert()
	if err != nil {
		return nil, err
	}
	return &l, nil
}
