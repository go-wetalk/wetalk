package service

import (
	"appsrv/model"
	"appsrv/schema"

	"github.com/go-pg/pg/v9"
	"github.com/xeonx/timeago"
)

type Topic struct{}

// ListWithRankByScore 综合评分进行排序的话题列表
func (Topic) ListWithRankByScore(db *pg.DB, input schema.TopicListInput) ([]schema.TopicListItem, error) {
	ts := []model.Topic{}
	err := db.Model(&ts).Order("topic.id DESC").Limit(input.Size).Column("topic.id", "topic.title", "topic.user_id", "topic.created").Relation("User").Select()

	out := []schema.TopicListItem{}
	for _, t := range ts {
		item := schema.TopicListItem{}
		item.ID = t.ID
		item.Title = t.Title
		item.Created = timeago.Chinese.Format(t.Created)

		if t.User != nil {
			item.User = &schema.User{
				ID:   t.UserID,
				Name: t.User.Name,
				Logo: t.User.Logo,
			}
		}

		if lastComment, err := t.LastComment(db); err == nil && lastComment != nil {
			item.LastComment = &schema.CommentBadge{
				ID:      lastComment.ID,
				TopicID: lastComment.TopicID,
				Created: timeago.Chinese.Format(lastComment.Created),
				User: &schema.User{
					ID:   lastComment.UserID,
					Name: lastComment.User.Name,
					Logo: lastComment.User.Logo,
				},
			}
		}

		out = append(out, item)
	}

	return out, err
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
