package model

import (
	"appsrv/pkg/db"
	"time"
)

const (
	AnnounceSlotShop uint8 = 1 // 试题
	AnnounceSlotH5   uint8 = 2 // H5页面
	AnnounceSlotText uint8 = 3 // 文本公告
)

// Announce 首页 Banner
type Announce struct {
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
