package service

import (
	"appsrv/model"
	"appsrv/schema"

	"github.com/go-pg/pg/v9"
)

type Topic struct{}

// ListWithRankByScore 综合评分进行排序的话题列表
func (Topic) ListWithRankByScore(db *pg.DB, input schema.TopicListInput) ([]model.Topic, error) {
	ts := []model.Topic{}
	err := db.Model(&ts).Order("id DESC").Limit(input.Size).Select()
	return ts, err
}

// Create 创建话题
func (Topic) Create(db *pg.DB, u model.User, input schema.TopicCreateInput) (*model.Topic, error) {
	t := model.Topic{}
	t.Title = input.Title
	t.Content = input.Content
	t.UserID = u.ID
	err := db.Insert(&t)
	return &t, err
}
