// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package service

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/oss"
	"github.com/google/wire"
)

// Injectors from wire.go:

func NewUserService() *User {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	user := &User{
		db:   pgDB,
		log:  logger,
		mc:   client,
		conf: serverConfig,
	}
	return user
}

func NewTopicService() *Topic {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	topic := &Topic{
		db:  pgDB,
		log: logger,
	}
	return topic
}

func NewNotificationService() *Notification {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	notification := &Notification{
		db:   pgDB,
		log:  logger,
		mc:   client,
		conf: serverConfig,
	}
	return notification
}

func NewCommentService() *Comment {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	comment := &Comment{
		db:   pgDB,
		log:  logger,
		mc:   client,
		conf: serverConfig,
	}
	return comment
}

// wire.go:

var ServiceSet = wire.NewSet(
	NewUserService,
	NewTopicService,
	NewNotificationService,
	NewCommentService,
)