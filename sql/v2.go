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
			&role{
				Key:   "user:ro",
				Name:  "用户检索",
				Intro: "允许检索用户资料",
			},
			&role{
				Key:   "user",
				Name:  "用户管理",
				Intro: "允许对用户信息进行编辑操作",
			},
		)

		return nil
	})
}

type user struct {
	ID        uint
	Name      string `pg:",notnull,unique"`
	Phone     string
	OpenID    string `pg:",unique"`
	AvatarURL string
	Gender    int  `pg:",default:1"`
	Coin      int  `pg:",default:0"`
	Remark    int8 `pg:",default:0"` // 账号来源标记
	Street    string
	City      string
	Province  string
	Region    string
	Email     string `pg:",unique"`
	Password  string

	db.TimeUpdate
}
