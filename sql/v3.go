package sql

import (
	"appsrv/pkg/db"
	"time"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func init() {
	Setup(3, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			&text{},
			&announce{},
		)

		db.Insert(
			&role{
				Key:   "text",
				Name:  "文本数据管理",
				Intro: "管理全局文本数据",
			},
			&role{
				Key:   "text:ro",
				Name:  "文本数据检索",
				Intro: "检索全局文本数据",
			},
		)

		db.Insert(
			&text{
				Name:     "服务协议",
				Slot:     1,
				SlotName: "service",
				Content:  "--terms OF service--",
			},
			&text{
				Name:     "隐私协议",
				Slot:     1,
				SlotName: "privacy",
				Content:  "--terms OF privacy--",
			},
			&text{
				Name:     "隐私协议",
				Slot:     3,
				SlotName: "guide",
				Content:  "--terms OF privacy--",
			},
		)

		return nil
	})
}

type text struct {
	ID       uint
	Name     string
	Slot     uint8  `pg:",unique:dedi"`
	SlotName string `pg:",unique:dedi"`
	Content  string

	db.TimeUpdate
}

type announce struct {
	ID        uint
	Name      string
	Show      *time.Time
	Hide      *time.Time
	Slot      uint8  // 根据 Slot 的值来决定行为
	SlotID    uint   `pg:",default:0"`
	SlotParam string // 文本参数就存这个字段，比如 wap 网页的地址
	Seq       int    `pg:",default:0"`

	db.LogoField
	db.TimeUpdate
}
