package weapp

import (
	"appsrv/pkg/config"
	"appsrv/pkg/redis"
	"sync"
	"time"

	"github.com/go-redis/cache/v7"
	"github.com/medivhzhan/weapp/v2"
)

const accessTokenKey = "cache:weapp:accesstoken"

// GetAccessToken 获取访问口令
func GetAccessToken() (string, error) {
	s := ""
	err := redis.Cache.Get(accessTokenKey, &s)
	if err != nil {
		for i := 0; i < 3; i++ { // 循环三次
			s, err = getAccessToken()
			if err == nil {
				break
			}
		}
	}

	return s, err
}

var l sync.Mutex

func getAccessToken() (s string, err error) {
	l.Lock()
	defer l.Unlock()

	res, err := weapp.GetAccessToken(config.Server.Weapp.AppID, config.Server.Weapp.Secret)
	if err != nil {
		return "", err
	}

	if err = res.GetResponseError(); err != nil {
		return "", err
	}

	s = res.AccessToken

	err = redis.Cache.Set(&cache.Item{
		Key:        accessTokenKey,
		Object:     s,
		Expiration: time.Duration(res.ExpiresIn) * time.Second,
	})

	return
}
