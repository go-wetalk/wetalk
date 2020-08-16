package redis

import (
	"appsrv/pkg/config"

	goredis "github.com/go-redis/redis/v7"
)

var rc *goredis.Client

func ProvideSingleton() *goredis.Client {
	if rc == nil {
		srvConf := config.ProvideSingleton()
		rc = goredis.NewClient(&goredis.Options{
			Addr:     srvConf.Redis.Addr,
			Password: srvConf.Redis.Password,
			DB:       srvConf.Redis.DB,
		})

		if err := rc.Ping().Err(); err != nil {
			panic(err)
		}
	}

	return rc
}
