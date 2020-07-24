package schema

type Comment struct {
	ID      uint
	TopicID uint
	Content string
	Created string

	User *User
}
