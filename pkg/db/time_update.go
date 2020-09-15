package db

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

type TimeUpdate struct {
	Created time.Time  `pg:",default:now(),index"`
	Updated time.Time  `pg:",default:now()"`
	Deleted *time.Time `pg:",soft_delete"`
}

var _ pg.BeforeUpdateHook = (*TimeUpdate)(nil)

func (tu *TimeUpdate) BeforeUpdate(ctx context.Context) (context.Context, error) {
	tu.Updated = time.Now()
	return ctx, nil
}
