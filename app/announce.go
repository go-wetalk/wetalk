//+build wireinject

package app

import (
	"appsrv/model"
	"appsrv/pkg"
	"appsrv/pkg/config"
	"appsrv/pkg/out"
	"appsrv/pkg/runtime"
	"net/http"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/google/wire"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"go.uber.org/zap"
)

type Announce struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}

func (v *Announce) RegisterRoute(m muxie.SubMux) {
	m.Handle("/announces", muxie.Methods().HandleFunc(http.MethodGet, v.AppList))
}

func NewAnnounceController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Announce), "*"),
		wire.Bind(new(runtime.Controller), new(*Announce)),
	)
	return nil
}

func (v Announce) List(w http.ResponseWriter, r *http.Request) {
	var as = []model.Announce{}
	_ = v.db.Model(&as).OrderExpr("seq DESC, id ASC").Select()
	muxie.Dispatch(w, muxie.JSON, &as)
}

func (v Announce) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Announce
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		v.log.Error("Announce.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	err = v.db.Insert(&in.Announce)
	if err != nil {
		v.log.Error("Announce.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (v Announce) Delete(w http.ResponseWriter, r *http.Request) {
	var a model.Announce
	err := v.db.Model(&a).Where("id = ?", muxie.GetParam(w, "announceID")).First()
	if err != nil {
		w.WriteHeader(404)
		return
	}

	_, err = v.db.Model(&a).WherePK().Delete()
	if err != nil {
		v.log.Error("Announce.Delete", zap.Error(err))
	}
	w.WriteHeader(204)
}

func (v Announce) AppList(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	as := []model.Announce{}
	_ = v.db.Model(&as).Where("(show IS NULL OR show < ?) AND (hide IS NULL OR hide > ?)", t, t).OrderExpr("seq DESC, id ASC").Select()

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
