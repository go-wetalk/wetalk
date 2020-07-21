package model

// SiteConfig 全站持久性配置
type SiteConfig struct {
	Key     string `pg:",unique,notnull"`
	Value   string
	Version uint `pg:",default:0"`
	Comment string
}
