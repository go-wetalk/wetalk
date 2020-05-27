package main

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/oss"
	"appsrv/pkg/redis"
	"appsrv/sql"
	"errors"
	"net/http"

	"github.com/kataras/muxie"
	"github.com/koding/multiconfig"
	"go.uber.org/zap"
)

func main() {
	initServices() // 初始化各项服务

	m := muxie.NewMux()
	m.PathCorrection = true
	m.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Add("Access-Control-Allow-Methods", "*")
			w.Header().Add("Access-Control-Allow-Headers", "Authorization,Content-Type")
			w.Header().Add("Access-Control-Max-Age", "600")
			w.Header().Add("Access-Control-Expose-Headers", "X-Refresh-Token")
			w.Header().Add("Vary", "Origin")
			if r.Method == http.MethodOptions {
				w.WriteHeader(200)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})

	initAdminServerV1(m.Of("/adm/v1"))
	initAppServerV1(m.Of("/app/v1"))

	if config.Server.Port == "" {
		config.Server.Port = ":8080"
	}

	bog.Info("server started", zap.String("addr", config.Server.Port))
	err := http.ListenAndServe(config.Server.Port, m)
	if err != nil {
		bog.Error("server error", zap.Error(err))
	}
}

func initServices() {
	bog.InitLog()

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
