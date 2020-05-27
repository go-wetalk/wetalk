package app

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"net/http"
	"time"

	"github.com/kataras/muxie"
	"go.uber.org/zap"
)

type AdminLog struct {
	ID        uint
	AdminID   uint
	AdminName string
	Event     string
	Intro     string
	IP        string
	UA        string
	Ref       string
	Created   time.Time `pg:",default:now()"`
}

func (AdminLog) LogEvent(r *http.Request, a *Admin, e, intro string) {
	db.DB.Insert(&AdminLog{
		AdminID:   a.ID,
		AdminName: a.Name,
		Event:     e,
		Intro:     intro,
		IP:        r.RemoteAddr,
		UA:        r.Header.Get("User-Agent"),
		Ref:       r.Referer(),
	})
}

func (AdminLog) List(w http.ResponseWriter, r *http.Request) {
	csr := r.URL.Query().Get("_csr")
	var ls = []AdminLog{}
	q := db.DB.Model(&ls).Order("admin_log.id desc").Limit(15)
	if csr != "0" {
		q = q.Where("admin_log.id < ?", csr)
	}
	err := q.Select()
	if err != nil {
		bog.Error("AdminLog.List", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, &ls)
}
