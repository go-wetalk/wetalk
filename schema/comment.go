package schema

import (
	"math"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Comment struct {
	CommentBadge

	Content string
}

type CommentBadge struct {
	ID      uint
	TopicID uint
	Created string

	User *User
}

type TopicCommentCreation struct {
	TopicID uint
	Content string
}

func (v TopicCommentCreation) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.TopicID, validation.Required.Error("无指定帖子")),
		validation.Field(&v.Content, validation.RuneLength(4, math.MaxUint32).Error("评论内容至少4个字")),
	)
}
