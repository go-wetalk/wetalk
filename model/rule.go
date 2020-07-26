package model

import "appsrv/pkg/db"

// Rule RBAC配置表
type Rule struct {
	ID          uint
	Host        string `pg:",unique:action,default:*"`
	Path        string `pg:",unique:action,default:*"`
	Method      string `pg:",unique:action,default:*"`
	AllowAnyone bool
	Authorized  []string `pg:",array"`
	Forbidden   []string `pg:",array"`

	db.TimeUpdate
}
