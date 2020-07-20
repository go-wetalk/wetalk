package model

import (
	"appsrv/pkg/db"
	"context"

	"github.com/go-pg/pg/v9"
)

type CoinLog struct {
	ID       uint
	UserID   uint
	Source   string // 变动因素
	SourceID uint
	Value    int
	Balance  int

	db.TimeUpdate
}

var _ pg.AfterInsertHook = (*CoinLog)(nil)

func (cl *CoinLog) AfterInsert(ctx context.Context) error {
	_, err := db.DB.Model(&User{ID: cl.UserID}).WherePK().Set("coin = ?", cl.Balance).Update()
	return err
}
