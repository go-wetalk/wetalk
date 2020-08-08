package sql

import (
	"appsrv/pkg/db"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func init() {
	Setup(7, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			&notification{},
			&message{},
		)
		return nil
	})
}

// 私信
type message struct {
	ID      uint
	SendID  uint  // 发送方UID
	RecvID  uint  // 接收方UID
	Kind    uint8 // 消息类型
	Content string

	db.TimeUpdate
}

// 系统通知
type notification struct {
	ID      uint
	RecvID  uint
	Content string
	HasRead bool `pg:",default:false"`

	db.TimeUpdate
}
