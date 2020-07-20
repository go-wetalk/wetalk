package model

import "appsrv/pkg/db"

const (
	TextSlotTerms    uint8 = 1 // 条款
	TextSlotAnnounce uint8 = 2 // 公告
	TextSlotNotice   uint8 = 3 // 提示
)

type Text struct {
	ID       uint
	Name     string
	Slot     uint8
	SlotName string `pg:",unique"`
	Content  string

	db.TimeUpdate
}
