package redis

import (
	gocache "github.com/go-redis/cache/v7"
	goredis "github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"
)

var cc *gocache.Codec

// ProvideCache provides cache wrapper for redis.
func ProvideCache(rc *goredis.Client) *gocache.Codec {
	if cc == nil {
		cc = &gocache.Codec{
			Redis: rc,
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		}
	}

	return cc
}
