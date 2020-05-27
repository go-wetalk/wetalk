package app

import (
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/now"
	"github.com/kataras/muxie"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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
func (t *Task) AvailableTo(u *User) bool {
	n := time.Now()
	if t.Begin != nil && t.Begin.After(n) {
		return false
	}
	if t.End != nil && t.End.Before(n) {
		return false
	}

	l := TaskLog{}
	err := db.DB.Model(&l).Where("task_id = ? AND user_id = ?", t.ID, u.ID).Order("id DESC").First()
	if err != nil {
		bog.Error("Task.AvailableTo", zap.Error(err), zap.Uint("UserID", u.ID), zap.Uint("TaskID", t.ID))
		return true
	}

	// 每日刷新的任务，按天计算，不按小时计算
	if t.Daily && l.Created.Before(now.BeginningOfDay()) {
		return true
	}

	if t.Cooling > 0 && l.Created.Add(time.Duration(t.Cooling)*time.Second).Before(n) {
		return true
	}

	count, _ := db.DB.Model((*TaskLog)(nil)).Where("task_id = ? AND user_id = ?", t.ID, u.ID).Count()
	if t.Times > 0 && int(t.Times) > count {
		return true
	}

	return false
}

func (t *Task) StatusText(u *User) string {
	switch t.Bonus {
	case BonusCoin:
		if t.Fulfilled(u) {
			if t.Got(u) {
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

func (t *Task) Got(u *User) bool {
	count, _ := db.DB.Model((*TaskLog)(nil)).
		Where("task_id = ? AND user_id = ?", t.ID, u.ID).
		Where("created BETWEEN ? AND ?", now.BeginningOfDay(), now.EndOfDay()).
		Count()
	return count > 0
}

func (t *Task) Confirm(u *User) (*TaskLog, error) {
	l := TaskLog{
		TaskID:    t.ID,
		UserID:    u.ID,
		Bonus:     t.Bonus,
		BonusNum:  t.BonusNum,
		Factor:    t.Factor,
		FactorNum: t.FactorNum,
	}
	err := db.DB.Insert(&l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// ******** 控制器逻辑

func (Task) List(w http.ResponseWriter, r *http.Request) {
	var ts = []Task{}
	_ = db.DB.Model(&ts).Order("id ASC").Select()
	muxie.Dispatch(w, muxie.JSON, &ts)
}

func (Task) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Task
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Task.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	err = db.DB.Insert(&in.Task)
	if err != nil {
		bog.Error("Task.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (Task) Update(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name    string
		Content string
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Task.Update", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	var t Task
	err = db.DB.Model(&t).Where("id = ?", muxie.GetParam(w, "TaskID")).First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	_, err = db.DB.Model(&t).WherePK().Set("name = ?", in.Name).Set("content = ?", in.Content).Update()
	if err != nil {
		bog.Error("Task.Update", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (Task) Delete(w http.ResponseWriter, r *http.Request) {
	m := Task{}
	err := db.DB.Model(&m).Where("id = ?", muxie.GetParam(w, "taskID")).First()
	if err == nil {
		_, err = db.DB.Model(&m).WherePK().Delete()
		if err != nil {
			bog.Error("Task.Delete", zap.Error(err))
			w.WriteHeader(500)
			return
		}
	}

	w.WriteHeader(204)
}

func (Task) AppList(w http.ResponseWriter, r *http.Request) {
	ts := []Task{}
	now := time.Now()
	count, err := db.DB.Model(&ts).
		Where("(\"begin\" IS NULL OR \"begin\" <= ?) AND (\"end\" IS NULL OR \"end\" >= ?)", now, now).
		Order("seq DESC").
		SelectAndCount()
	if err != nil {
		bog.Error("Task.AppList", zap.Error(err))
	}

	var u User
	auth.GetUser(r, &u)

	type outItem struct {
		ID         uint
		Name       string
		Intro      string
		Button     bool // 按钮形式就表示已完成未领取
		StatusText string
	}

	var out = struct {
		Count int
		Page  int
		Data  []outItem
	}{
		Count: count,
		Page:  1,
		Data:  []outItem{},
	}
	for _, t := range ts {
		if !t.AvailableTo(&u) {
			continue
		}

		i := outItem{
			ID:         t.ID,
			Name:       t.Name,
			Intro:      t.Intro,
			Button:     t.Fulfilled(&u) && !t.Got(&u),
			StatusText: t.StatusText(&u),
		}
		out.Data = append(out.Data, i)
	}
	muxie.Dispatch(w, muxie.JSON, &out)
}

func (Task) AppTaskLogCreate(w http.ResponseWriter, r *http.Request) {
	var u User
	err := auth.GetUser(r, &u)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var t Task
	kword := muxie.GetParam(w, "taskID")
	err = db.DB.Model(&t).Where("id = ? or factor = ?", cast.ToInt(kword), kword).Order("id DESC").First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if !t.AvailableTo(&u) {
		w.WriteHeader(404)
		return
	}

	if !t.Fulfilled(&u) {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	if t.Got(&u) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	_, err = t.Confirm(&u)
	if err != nil {
		bog.Error("Task.AppTaskLogCreate", zap.Error(err), zap.Uint("TaskID", t.ID), zap.Uint("UserID", u.ID))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}
