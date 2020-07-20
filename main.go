package main

import (
	"appsrv/cmd"
	"appsrv/pkg/bog"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/oss"
	"appsrv/pkg/redis"
	"appsrv/sql"
	"errors"

	"github.com/koding/multiconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	cobra.OnInitialize(initServices)

	if err := cmd.RootCommand.Execute(); err != nil {
		bog.Fatal("Basis.Execute", zap.Error(err))
	}
}

func initServices() {
	l := multiconfig.MultiLoader(
		&multiconfig.EnvironmentLoader{},
		&multiconfig.TOMLLoader{
			Path: "local.toml",
		},
	)
	err := l.Load(&config.Server)
	if !errors.Is(err, multiconfig.ErrFileNotFound) {
		checkErr(err)
	}

	err = redis.InitRedis(config.Server.Redis)
	checkErr(err)
	bog.Info("redis connected", zap.String("addr", config.Server.Redis.Addr))

	redis.InitCache()

	err = db.InitConn(config.Server.DB, config.Server.Debug)
	checkErr(err)
	bog.Info("database connected", zap.String("addr", config.Server.DB.Addr))

	sql.Run()

	err = oss.InitOss(config.Server.Oss)
	checkErr(err)
	bog.Info("minio connected", zap.String("addr", config.Server.Oss.Endpoint))
}

func checkErr(err error) {
	if err != nil {
		bog.Fatal("boot error", zap.Error(err))
	}
}
