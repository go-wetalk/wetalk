package model

import (
	"appsrv/pkg/db"

	"github.com/go-pg/pg/v10"
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

// var _ pg.AfterInsertHook = (*User)(nil)

// func (u User) AfterInsert(c context.Context) error {
// 	n := Notification{
// 		RecvID:  u.ID,
// 		Content: "欢迎你，优秀的第 " + fmt.Sprintf("%d", u.ID) + " 号会员。祝你玩的开心，能有更多收获与积累。",
// 	}
// 	return db.DB.Insert(&n)
// }

func (u User) UnreadNotify(db *pg.DB) (count int) {
	count, _ = db.Model((*Notification)(nil)).Where("recv_id = ? and has_read = ?", u.ID, false).Count()
	return
}
