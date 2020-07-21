package service

import (
	"appsrv/schema"

	"github.com/go-pg/pg/v9"
)

type Topic struct{}

// ListWithRankByScore 综合评分进行排序的话题列表
func (Topic) ListWithRankByScore(db *pg.DB, input schema.TopicListInput) ([]Topic, error) {
	ts := []Topic{}
	err := db.Model(&ts).Order("id DESC").Limit(input.Size).Select()
	return ts, err
}
