package route

import (
	"appsrv/app"
	"net/http"

	"github.com/kataras/muxie"
)

// SetupAppServerV1 setup user-end routes.
func SetupAppServerV1(v1 muxie.SubMux) {
	v1.Handle("/users", muxie.Methods().HandleFunc(http.MethodPost, app.User{}.SignUp))
	v1.Handle("/tokens", muxie.Methods().HandleFunc(http.MethodPost, app.User{}.Login))

	v1.Handle("/announces", muxie.Methods().HandleFunc(http.MethodGet, app.Announce{}.AppList))

	v1.Handle("/texts/:textID", muxie.Methods().HandleFunc(http.MethodGet, app.Text{}.AppView))

	v1.Handle("/status", muxie.Methods().HandleFunc(http.MethodGet, app.User{}.AppStatus))

	v1.Handle("/topics", muxie.Methods().
		HandleFunc(http.MethodGet, app.Topic{}.List).
		HandleFunc(http.MethodPost, app.Topic{}.Create))

	v1.Handle("/topics/:topicID", muxie.Methods().
		HandleFunc(http.MethodGet, app.Topic{}.Find))

	v1.Handle("/comments", muxie.Methods().
		HandleFunc(http.MethodGet, app.Comment.ListByFilter).
		HandleFunc(http.MethodPost, app.Comment.CreateTopicComment))

	v1.Handle("/tasks", muxie.Methods().
		HandleFunc(http.MethodGet, app.Task{}.AppList))
	v1.Handle("/tasks/:taskID/bonus", muxie.Methods().
		HandleFunc(http.MethodPost, app.Task{}.AppTaskLogCreate))

	v1.Handle("/users/:name", muxie.Methods().HandleFunc(http.MethodGet, app.User{}.ViewUserDetail))

	v1.Handle("/notifications", muxie.Methods().HandleFunc(http.MethodGet, app.Notification.List))
	v1.Handle("/notifications/:notificationID", muxie.Methods().HandleFunc(http.MethodDelete, app.Notification.MarkRead))
}
