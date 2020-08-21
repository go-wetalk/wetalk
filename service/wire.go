//+build wireinject

package service

import (
	"appsrv/pkg"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	NewUserService,
	NewTopicService,
	NewNotificationService,
	NewCommentService,
)

func NewUserService() *User {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(User), "*"),
	)
	return nil
}

func NewTopicService() *Topic {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Topic), "*"),
	)
	return nil
}

func NewNotificationService() *Notification {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Notification), "*"),
	)
	return nil
}

func NewCommentService() *Comment {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Comment), "*"),
	)
	return nil
}
