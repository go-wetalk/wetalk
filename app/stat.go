package app

import (
	"appsrv/pkg/db"
	"net/http"
	"time"

	"github.com/kataras/muxie"
)

// Stat 数据统计、报表相关逻辑
type Stat struct{}

func (Stat) Summary(w http.ResponseWriter, r *http.Request) {
	var out struct {
		User           int // 用户总数
		LastDayNewUser int
		NewUser        int
		LastDayOrder   int
		Order          int
		User7Days      []int
		Order7Days     []int
	}
	out.User, _ = db.DB.Model((*User)(nil)).Count()

	t := time.Now().AddDate(0, 0, -1)
	yesterdayBegin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	yesterdayEnd := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())

	out.LastDayNewUser, _ = db.DB.Model((*User)(nil)).Where("created BETWEEN ? AND ?", yesterdayBegin, yesterdayEnd).Count()

	t = time.Now()
	todayBegin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	todayEnd := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())

	out.NewUser, _ = db.DB.Model((*User)(nil)).Where("created BETWEEN ? AND ?", todayBegin, todayEnd).Count()

	for i := 6; i >= 0; i-- {
		t = time.Now().AddDate(0, 0, -i)
		begin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())

		cnt, _ := db.DB.Model((*User)(nil)).Where("created BETWEEN ? AND ?", begin, end).Count()
		out.User7Days = append(out.User7Days, cnt)
	}

	muxie.Dispatch(w, muxie.JSON, &out)
}
