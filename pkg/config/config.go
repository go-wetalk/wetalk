package config

type ServerConfig struct {
	Debug bool
	Port  string
	DB    DBConfig
	Redis RedisConfig
	Oss   OssConfig
	Weapp WeappConfig
	Qapp  QappConfig
	Auth  AuthConfig
}

var Server ServerConfig
