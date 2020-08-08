package schema

type Notification struct {
	ID      uint
	RecvID  uint
	Content string
	HasRead bool
	Created string
}
