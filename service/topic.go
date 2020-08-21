package service

import (
	"appsrv/model"
	"appsrv/pkg/out"
	"appsrv/schema"
	"bytes"
	"strings"

	"github.com/88250/lute"
	"github.com/go-pg/pg/v9"
	"github.com/xeonx/timeago"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"go.uber.org/zap"
)

// Topic 主题相关DB操作
type Topic struct {
	db  *pg.DB
	log *zap.Logger
}

// ListWithRankByScore 综合评分进行排序的主题列表
func (v *Topic) ListWithRankByScore(input schema.TopicListInput) (*schema.Pagination, error) {
	pag := schema.Pagination{
		PerPage: input.Size,
	}

	ts := []model.Topic{}
	q := v.db.Model(&ts).Column("topic.id", "topic.title", "topic.user_id", "topic.created", "topic.tags").
		Relation("User").
		Order("topic.id DESC").
		Offset((input.Page - 1) * input.Size).Limit(input.Size)
	if input.Tag != "" {
		q.Where("? = ANY(topic.tags)", input.Tag)
	}
	count, err := q.SelectAndCount()
	if err != nil {
		v.log.Error("topic.ListWithRankByScore", zap.Error(err))
		return nil, out.Err500
	}

	pag.RowCount = count

	data := []schema.TopicListItem{}
	for _, t := range ts {
		item := schema.TopicListItem{}
		item.ID = t.ID
		item.Title = t.Title
		item.Tags = t.Tags
		item.Created = timeago.Chinese.Format(t.Created)

		if t.User != nil {
			item.User = &schema.User{
				ID:   t.UserID,
				Name: t.User.Name,
				Logo: t.User.LogoLink(),
			}
		}

		if lastComment, err := t.LastComment(v.db); err == nil && lastComment != nil {
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

		data = append(data, item)
	}

	pag.Data = data

	return &pag, err
}

// Create 创建主题
func (v *Topic) Create(u model.User, input schema.TopicCreateInput) (*model.Topic, error) {
	lu := lute.New()
	t := model.Topic{}
	t.Title = input.Title
	t.Content = lu.FormatStr("", input.Content)
	t.UserID = u.ID

	for _, tag := range input.Tags {
		if tag = strings.TrimSpace(tag); tag != "" {
			t.Tags = append(t.Tags, tag)
		}
	}

	err := v.db.Insert(&t)
	return &t, err
}

// FindByID 查找主题
func (v *Topic) FindByID(id uint) (*schema.Topic, error) {
	t := model.Topic{}
	err := v.db.Model(&t).Relation("User").Where("topic.id = ?", id).First()
	if err != nil {
		return nil, out.Err500
	}

	item := schema.TopicListItem{}
	item.ID = t.ID
	item.Title = t.Title
	item.Tags = t.Tags
	item.Created = timeago.Chinese.Format(t.Created)

	if t.User != nil {
		item.User = &schema.User{
			ID:   t.UserID,
			Name: t.User.Name,
			Logo: t.User.LogoLink(),
		}
	}

	if lastComment, err := t.LastComment(v.db); err == nil && lastComment != nil {
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

	to := schema.Topic{
		TopicListItem: item,
		Content:       t.Content,
	}

	gm := goldmark.New(goldmark.WithExtensions(extension.GFM))

	var b bytes.Buffer
	err = gm.Convert([]byte(t.Content), &b)
	if err != nil {
		v.log.Error("goldmark.Convert", zap.Error(err))
		return nil, out.Err500
	}

	to.Content = b.String()

	return &to, nil
}
