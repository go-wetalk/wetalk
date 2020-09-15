package service

import (
	"appsrv/model"
	"appsrv/pkg/config"
	"appsrv/pkg/out"
	"appsrv/schema"

	"github.com/go-pg/pg/v10"
	"github.com/minio/minio-go/v6"
	"github.com/xeonx/timeago"
	"go.uber.org/zap"
)

// Notification 通知相关DB操作
type Notification struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}

func (v *Notification) FindForUser(u *model.User, input schema.Paginate) (ret schema.Pagination, err error) {
	ret.PerPage = input.Size

	ns := []model.Notification{}
	ret.RowCount, err = v.db.Model(&ns).
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

func (v *Notification) MarkAsRead(u *model.User, notifyID uint) (err error) {
	_, err = v.db.Model((*model.Notification)(nil)).
		Where("id = ? and recv_id = ?", notifyID, u.ID).
		Set("has_read = ?", true).Update()
	return
}
