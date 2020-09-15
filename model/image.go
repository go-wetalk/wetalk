package model

import (
	"appsrv/pkg/db"
	"context"

	"github.com/go-pg/pg/v10"
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
	// return oss.Server + "/" + oss.Bucket + "/" + i.Path
	return ""
}
