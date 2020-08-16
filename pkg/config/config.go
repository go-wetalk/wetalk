package config

import (
	"errors"
	"log"

	"github.com/koding/multiconfig"
)

// ServerConfig provides runtime settings to server.
type ServerConfig struct {
	Debug bool
	Port  string
	DB    DBConfig
	Redis RedisConfig
	Oss   OssConfig
	Weapp WeappConfig
	Qapp  QappConfig
	Auth  AuthConfig

	HCaptcha struct {
		Enabled bool
		Secret  string
	}
}

var srvConf *ServerConfig

// ProvideSingleton provides ServerConfig in singleton instance.
func ProvideSingleton() *ServerConfig {
	if srvConf == nil {
		v := ServerConfig{}
		l := multiconfig.MultiLoader(
			&multiconfig.EnvironmentLoader{},
			&multiconfig.TOMLLoader{
				Path: "local.toml",
			},
		)
		err := l.Load(&v)
		if err != nil && !errors.Is(err, multiconfig.ErrFileNotFound) {
			log.Fatal(err)
		}

		srvConf = &v
	}

	return srvConf
}
