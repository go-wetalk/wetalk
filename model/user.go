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
	Name     string `pg:",notnull,unique"`
	Phone    string `pg:",unique"`
	Email    string `pg:",unique"`
	Password string `json:"-"`
	Gender   int    `pg:",default:1"`
	Coin     int    `pg:",default:0"`
	Street   string
	City     string
	Province string
	Country  string
	Sign     string
	RoleKeys []string `pg:",array"`

	db.LogoField
	db.TimeUpdate
}
