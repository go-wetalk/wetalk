package app

import (
	"appsrv/model"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"appsrv/pkg/out"
	"net/http"
	"time"

	"github.com/kataras/muxie"
	"go.uber.org/zap"
)

type Announce struct{}

func (Announce) List(w http.ResponseWriter, r *http.Request) {
	var as = []model.Announce{}
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
	var a model.Announce
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
	as := []model.Announce{}
	_ = db.DB.Model(&as).Where("(show IS NULL OR show < ?) AND (hide IS NULL OR hide > ?)", t, t).OrderExpr("seq DESC, id ASC").Select()

	var raw = []struct {
		ID        uint
		Name      string
		Slot      uint8
		SlotID    uint
		SlotParam string
		Seq       int
		Logo      string
	}{}
	for _, a := range as {
		raw = append(raw, struct {
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
	muxie.Dispatch(w, muxie.JSON, out.Data(raw))
}
