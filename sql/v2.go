package sql

import (
	"appsrv/pkg/db"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func init() {
	Setup(2, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			&user{},
		)

		db.Insert(
			&rule{
				Path:        "/users",
				Method:      "{POST}",
				AllowAnyone: true, // 用户注册允许所有人请求
			},
			&rule{
				Path:        "/tokens",
				Method:      "{POST}",
				AllowAnyone: true, // 用户登录允许所有人请求
			},
			&rule{
				Path:       "/users/*",
				Method:     "{POST,PUT,DELETE}",
				Authorized: []string{"*"},
			},
			&rule{
				Path:       "/tokens",
				Method:     "{PUT,DELETE}",
				Authorized: []string{"*"},
			},
			&rule{
				Path:       "/status",
				Authorized: []string{"*"},
			},
		)

		return nil
	})
}

type user struct {
	ID       uint
	Name     string `pg:",notnull,unique"`
	Phone    string
	OpenID   string `pg:",unique"`
	Gender   int    `pg:",default:1"`
	Coin     int    `pg:",default:0"`
	Remark   int8   `pg:",default:0"` // 账号来源标记
	Street   string
	City     string
	Province string
	Region   string
	Email    string `pg:",unique"`
	Password string
	Sign     string
	RoleKeys []string `pg:",array"`

	db.LogoField
	db.TimeUpdate
}
