package app

import (
	"appsrv/pkg/config"

	"github.com/go-pg/pg/v10"
	"github.com/minio/minio-go/v6"
	"go.uber.org/zap"
)

type Role struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}
