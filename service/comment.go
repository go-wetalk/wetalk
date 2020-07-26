package service

import (
	"appsrv/model"
	"appsrv/pkg/errors"
	"appsrv/schema"

	"github.com/go-pg/pg/v9"
)

type Comment struct{}

func (Comment) CreateTopicComment(db *pg.DB, u model.User, input schema.TopicCommentCreation) (*model.Comment, error) {
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
