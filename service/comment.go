package service

import (
	"appsrv/model"
	"appsrv/pkg/errors"
	"appsrv/schema"

	"github.com/go-pg/pg/v9"
	"github.com/xeonx/timeago"
)

var Comment = new(comment)

type comment struct{}

func (comment) CreateTopicComment(db *pg.DB, u model.User, input schema.TopicCommentCreation) (*model.Comment, error) {
	t := model.Topic{}
	err := db.Model(&t).Where("id = ?", input.TopicID).First()
	if err != nil {
		return nil, errors.ErrNotFound
	}

	com := model.Comment{}
	com.TopicID = input.TopicID
	com.UserID = u.ID
	com.Content = input.Content
	err = db.Insert(&com)
	if err != nil {
		return nil, errors.New(500, "服务器爆炸啦")
	}

	return &com, nil
}

func (comment) FindByFilterInput(db *pg.DB, input schema.CommentFilter) ([]schema.Comment, error) {
	cs := []model.Comment{}
	err := db.Model(&cs).Relation("User").Where("comment.topic_id = ?", input.TopicID).Order("comment.updated DESC").Offset((input.Page - 1) * 20).Limit(20).Select()
	if err != nil {
		return nil, errors.Err500
	}

	out := []schema.Comment{}
	for _, com := range cs {
		out = append(out, schema.Comment{
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

	return out, nil
}
