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

// Comment 评论相关DB操作
type Comment struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}

func (v *Comment) CreateTopicComment(u model.User, input schema.TopicCommentCreation) (*model.Comment, error) {
	t := model.Topic{}
	err := v.db.Model(&t).Where("id = ?", input.TopicID).First()
	if err != nil {
		return nil, out.ErrNotFound
	}

	com := model.Comment{}
	com.TopicID = input.TopicID
	com.UserID = u.ID
	com.Content = input.Content
	_, err = v.db.Model(&com).Insert()
	if err != nil {
		return nil, out.Err500
	}

	return &com, nil
}

func (v *Comment) FindByFilterInput(input schema.CommentFilter) (*schema.Pagination, error) {
	raw := schema.Pagination{
		PerPage: input.Size,
	}

	cs := []model.Comment{}
	count, err := v.db.Model(&cs).Relation("User").
		Where("comment.topic_id = ?", input.TopicID).
		Order("comment.id DESC").
		Offset((input.Page - 1) * input.Size).Limit(input.Size).
		SelectAndCount()
	if err != nil {
		return nil, out.Err500
	}

	raw.RowCount = count

	data := []schema.Comment{}
	for _, com := range cs {
		data = append(data, schema.Comment{
			CommentBadge: schema.CommentBadge{
				ID:      com.ID,
				TopicID: com.TopicID,
				Created: timeago.Chinese.Format(com.Created),
				User: &schema.User{
					ID:   com.UserID,
					Name: com.User.Name,
					Logo: com.User.LogoLink(),
				},
			},
			Content: com.Content,
		})
	}

	raw.Data = data

	return &raw, nil
}
