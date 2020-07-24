package schema

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Topic struct {
	TopicListItem
	Content string
}

type TopicListItem struct {
	ID          uint
	Title       string
	Created     string
	User        *User
	LastComment *Comment
}

type TopicListInput struct {
	Page uint
	Size int
}

type TopicCreateInput struct {
	Title   string
	Content string
	Tags    []string
}

func (v TopicCreateInput) Validate() error {
	return validation.ValidateStruct(
		&v,
		validation.Field(&v.Title, validation.Required.Error("标题不能为空"), validation.RuneLength(1, 50).Error("字符数量在1到50之间")),
		validation.Field(&v.Content, validation.Required.Error("请输入正文")),
	)
}
