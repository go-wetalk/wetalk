package model

import (
	"appsrv/pkg/db"
)

const (
	// UserRemarkWechat 微信端用户标记
	UserRemarkWechat = 1
	// UserRemarkQQ QQ端用户标记
	UserRemarkQQ = 2
)

type User struct {
	ID       uint
	Name     string `pg:",notnull"`
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

	db.LogoField
	db.TimeUpdate
}
