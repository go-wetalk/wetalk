//+build wireinject

package service

import (
	"appsrv/pkg"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	NewUserService,
)

func NewUserService() *User {
	wire.Build(
		pkg.ApplicationSet,
		wire.Struct(new(User), "*"),
	)
	return nil
}
