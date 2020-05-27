package app

import (
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"appsrv/pkg/oss"
	"appsrv/pkg/weapp"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// Image 图片资源
type Image struct {
	ID    uint
	Intro string `pg:",notnull"`
	Path  string `pg:",notnull"`

	db.TimeUpdate

	Link string `pg:"-"`
}

var _ pg.AfterSelectHook = (*Image)(nil)

func (i *Image) AfterSelect(c context.Context) error {
	i.Link = i.ImageLink()
	return nil
}

var _ pg.AfterScanHook = (*Image)(nil)

func (i *Image) AfterScan(c context.Context) error {
	i.Link = i.ImageLink()
	return nil
}

func (i *Image) ImageLink() string {
	return oss.Server + "/" + oss.Bucket + "/" + i.Path
}

func (Image) PreSign(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string
		Size int64
		Type string
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Image.PreSign", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	policy := minio.NewPostPolicy()
	_ = policy.SetBucket(oss.Bucket)
	_ = policy.SetContentLengthRange(in.Size, in.Size)
	_ = policy.SetContentType(in.Type)
	_ = policy.SetExpires(time.Now().Add(15 * time.Minute))
	_ = policy.SetKey(time.Now().Format("20060102") + "/" + in.Name)
	u, form, err := oss.Min.PresignedPostPolicy(policy)
	if err != nil {
		bog.Error("Image.PreSign", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	var out struct {
		URL  string
		Form map[string]string
	}
	out.URL = u.String()
	out.Form = form

	muxie.Dispatch(w, muxie.JSON, &out)
}

func (Image) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Intro string
		Path  string
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Image.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	err = db.DB.Insert(&Image{
		Intro: in.Intro,
		Path:  in.Path,
	})
	if err != nil {
		bog.Error("Image.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (Image) List(w http.ResponseWriter, r *http.Request) {
	var is = []Image{}
	_ = db.DB.Model(&is).Column("id", "intro", "path").Select()
	muxie.Dispatch(w, muxie.JSON, &is)
}

// AppWxaCodeView 查看小程序码，没有就生成
func (Image) AppWxaCodeView(w http.ResponseWriter, r *http.Request) {
	var u User
	_ = auth.GetUser(r, &u)

	path := fmt.Sprintf("%s?uid=%d", "pages/index/index", u.ID)
	targetID := cast.ToInt(r.URL.Query().Get("t"))
	switch targetID {
	case 1:
		id := cast.ToUint(r.URL.Query().Get("id"))
		if id == 0 {
			w.WriteHeader(400)
			return
		}
		path = fmt.Sprintf("%s?id=%d&uid=%d", "pages/index/detail", id, u.ID)
	case 2:
		id := cast.ToUint(r.URL.Query().Get("id"))
		if id == 0 {
			w.WriteHeader(400)
			return
		}
		path = fmt.Sprintf("%s?id=%d&uid=%d", "pages/market/detail", id, u.ID)
	}

	img := Image{}
	err := db.DB.Model(&img).Where("intro = ?", path).Order("id DESC").First()
	if err == nil {
		muxie.Dispatch(w, muxie.JSON, &struct{ Link string }{Link: img.ImageLink()})
		return
	}

	filepath, err := weapp.GenerateWxaCode(path)
	if err != nil {
		bog.Error("weapp.GenerateWxaCode", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	img.Intro = path
	img.Path = filepath
	err = db.DB.Insert(&img)
	if err != nil {
		bog.Error("Image.AppWxaCodeView", zap.Error(err), zap.String("Path", path))
		w.WriteHeader(500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, &struct{ Link string }{Link: img.ImageLink()})
}
