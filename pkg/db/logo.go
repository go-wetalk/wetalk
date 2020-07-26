package db

import (
	"appsrv/pkg/oss"
	"context"
	"strings"

	"github.com/go-pg/pg/v9"
)

type LogoField struct {
	LogoID   uint `pg:",default:0"`
	LogoPath string

	Logo string `pg:"-"`
}

var _ pg.AfterSelectHook = (*LogoField)(nil)

func (l *LogoField) AfterSelect(c context.Context) error {
	l.Logo = l.LogoLink()
	return nil
}

var _ pg.AfterScanHook = (*LogoField)(nil)

func (l *LogoField) AfterScan(c context.Context) error {
	l.Logo = l.LogoLink()
	return nil
}

func (l LogoField) LogoLink() (s string) {
	if l.LogoPath != "" {
		if strings.HasPrefix(l.LogoPath, "https://") || strings.HasPrefix(l.LogoPath, "http://") {
			s = l.LogoPath
		} else {
			s = oss.Server + "/" + l.LogoPath
		}
	}
	return
}
