//+build wireinject

package app

import (
	"appsrv/pkg"
	"appsrv/pkg/runtime"
	"appsrv/service"

	"github.com/google/wire"
)

func NewUserController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		service.ServiceSet,
		wire.Struct(new(User), "*"),
		wire.Bind(new(runtime.Controller), new(*User)),
	)
	return nil
}

func NewTopicController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Topic), "*"),
		wire.Bind(new(runtime.Controller), new(*Topic)),
	)
	return nil
}

func NewTextController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Text), "*"),
		wire.Bind(new(runtime.Controller), new(*Text)),
	)
	return nil
}

func NewTaskController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Task), "*"),
		wire.Bind(new(runtime.Controller), new(*Task)),
	)
	return nil
}

func NewNotificationController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Notification), "*"),
		wire.Bind(new(runtime.Controller), new(*Notification)),
	)
	return nil
}

func NewCommentController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Comment), "*"),
		wire.Bind(new(runtime.Controller), new(*Comment)),
	)
	return nil
}

func NewAnnounceController() runtime.Controller {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(Announce), "*"),
		wire.Bind(new(runtime.Controller), new(*Announce)),
	)
	return nil
}
