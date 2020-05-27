package redis

import (
	gocache "github.com/go-redis/cache/v7"
	"github.com/vmihailenco/msgpack/v4"
)

var Cache *gocache.Codec

func InitCache() {
	Cache = &gocache.Codec{
		Redis: Redis,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}
