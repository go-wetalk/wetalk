package app

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"net/http"
	"time"

	"github.com/kataras/muxie"
	"go.uber.org/zap"
)

const (
	AnnounceSlotShop uint8 = 1 // 试题
	AnnounceSlotH5   uint8 = 2 // H5页面
	AnnounceSlotText uint8 = 3 // 文本公告
)

// Announce 首页 Banner
type Announce struct {
	ID        uint
	Name      string
	Show      *time.Time
	Hide      *time.Time
	Slot      uint8  // 根据 Slot 的值来决定行为
	SlotID    uint   `pg:",default:0"`
	SlotParam string // 文本参数就存这个字段，比如 wap 网页的地址
	Seq       int    `pg:",default:0"`

	db.LogoField
	db.TimeUpdate
}

// ******** 控制器逻辑

func (Announce) List(w http.ResponseWriter, r *http.Request) {
	var as = []Announce{}
	_ = db.DB.Model(&as).OrderExpr("seq DESC, id ASC").Select()
	muxie.Dispatch(w, muxie.JSON, &as)
}

func (Announce) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Announce
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Announce.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	err = db.DB.Insert(&in.Announce)
	if err != nil {
		bog.Error("Announce.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (Announce) Delete(w http.ResponseWriter, r *http.Request) {
	var a Announce
	err := db.DB.Model(&a).Where("id = ?", muxie.GetParam(w, "announceID")).First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	_, err = db.DB.Model(&a).WherePK().Delete()
	if err != nil {
		bog.Error("Announce.Delete", zap.Error(err))
	}
	w.WriteHeader(204)
}

func (Announce) AppList(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	as := []Announce{}
	_ = db.DB.Model(&as).Where("(show IS NULL OR show < ?) AND (hide IS NULL OR hide > ?)", t, t).OrderExpr("seq DESC, id ASC").Select()

	var out = []struct {
		ID        uint
		Name      string
		Slot      uint8
		SlotID    uint
		SlotParam string
		Seq       int
		Logo      string
	}{}
	for _, a := range as {
		out = append(out, struct {
			ID        uint
			Name      string
			Slot      uint8
			SlotID    uint
			SlotParam string
			Seq       int
			Logo      string
		}{
			ID:        a.ID,
			Name:      a.Name,
			Slot:      a.Slot,
			SlotID:    a.SlotID,
			SlotParam: a.SlotParam,
			Seq:       a.Seq,
			Logo:      a.LogoLink(),
		})
	}
	muxie.Dispatch(w, muxie.JSON, &out)
}
