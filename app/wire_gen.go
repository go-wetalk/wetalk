// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package app

import (
	"appsrv/pkg/bog"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/oss"
	"appsrv/pkg/runtime"
	"appsrv/service"
)

// Injectors from wire.go:

func NewUserController() runtime.Controller {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	user := service.NewUserService()
	appUser := &User{
		db:          pgDB,
		log:         logger,
		mc:          client,
		conf:        serverConfig,
		userService: user,
	}
	return appUser
}

func NewTopicController() runtime.Controller {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	topic := service.NewTopicService()
	appTopic := &Topic{
		db:           pgDB,
		log:          logger,
		mc:           client,
		conf:         serverConfig,
		topicService: topic,
	}
	return appTopic
}

func NewTextController() runtime.Controller {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	text := &Text{
		db:   pgDB,
		log:  logger,
		mc:   client,
		conf: serverConfig,
	}
	return text
}

func NewTaskController() runtime.Controller {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	task := &Task{
		db:   pgDB,
		log:  logger,
		mc:   client,
		conf: serverConfig,
	}
	return task
}

func NewNotificationController() runtime.Controller {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	notification := service.NewNotificationService()
	appNotification := &Notification{
		db:                  pgDB,
		log:                 logger,
		mc:                  client,
		conf:                serverConfig,
		notificationService: notification,
	}
	return appNotification
}

func NewCommentController() runtime.Controller {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	comment := service.NewCommentService()
	appComment := &Comment{
		db:             pgDB,
		log:            logger,
		mc:             client,
		conf:           serverConfig,
		commentService: comment,
	}
	return appComment
}

func NewAnnounceController() runtime.Controller {
	pgDB := db.ProvideSingleton()
	logger := bog.ProvideSingleton()
	client := oss.ProvideSingleton()
	serverConfig := config.ProvideSingleton()
	announce := &Announce{
		db:   pgDB,
		log:  logger,
		mc:   client,
		conf: serverConfig,
	}
	return announce
}
