package redis

import (
	"appsrv/pkg/config"

	goredis "github.com/go-redis/redis/v7"
)

var Redis *goredis.Client

func InitRedis(c config.RedisConfig) error {
	Redis = goredis.NewClient(&goredis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})

	return Redis.Ping().Err()
}
