package pkg

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/oss"
	"appsrv/pkg/redis"

	"github.com/google/wire"
)

var ApplicationSet = wire.NewSet(
	config.ProvideSingleton,
	bog.ProvideSingleton,
	db.ProvideSingleton,
	oss.ProvideSingleton,
	redis.ProvideSingleton,
	redis.ProvideCache,
)
