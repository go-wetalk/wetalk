package service

import (
	"appsrv/model"
	"appsrv/pkg/out"
	"appsrv/schema"

	"github.com/go-pg/pg/v9"
	"github.com/xeonx/timeago"
)

var Notification = &notification{}

type notification struct{}

func (notification) FindForUser(db *pg.DB, u *model.User, input schema.Paginate) (ret schema.Pagination, err error) {
	ret.PerPage = input.Size

	ns := []model.Notification{}
	ret.RowCount, err = db.Model(&ns).
		Where("recv_id = ?", u.ID).Order("notification.id desc").
		Offset(input.Offset()).Limit(input.Size).
		SelectAndCount()
	if err != nil {
		return ret, out.Err500
	}

	out := []schema.Notification{}
	for _, n := range ns {
		out = append(out, schema.Notification{
			ID:      n.ID,
			RecvID:  n.RecvID,
			Content: n.Content,
			HasRead: n.HasRead,
			Created: timeago.Chinese.Format(n.Created),
		})
	}

	ret.Data = out
	return
}

func (notification) MarkAsRead(db *pg.DB, u *model.User, notifyID uint) (err error) {
	_, err = db.Model((*model.Notification)(nil)).
		Where("id = ? and recv_id = ?", notifyID, u.ID).
		Set("has_read = ?", true).Update()
	return
}
