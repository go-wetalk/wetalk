package model

import (
	"appsrv/pkg/db"

	"github.com/go-pg/pg/v10"
)

type Topic struct {
	ID      uint
	UserID  uint
	Title   string
	Content string
	Tags    []string `pg:",array"`

	db.TimeUpdate

	User     *User
	Comments []Comment
}

// LastComment 获取最后一条回复
func (t *Topic) LastComment(db *pg.DB) (*Comment, error) {
	c := Comment{}
	err := db.Model(&c).Where("topic_id = ?", t.ID).Order("id desc").Relation("User").First()
	return &c, err
}
