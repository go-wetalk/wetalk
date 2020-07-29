package service

import (
	"fmt"
	"net/url"
)

// SiteConfig 系统配置
var SiteConfig = new(siteConfig)

type siteConfig struct {
	Domain string
	Name   string
}

func (sc siteConfig) GenURL(path string, query url.Values) string {
	return fmt.Sprintf("https://%s%s%s", sc.Domain, path, query.Encode())
}
