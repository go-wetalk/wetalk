package route

import (
	"appsrv/app"
	"appsrv/pkg/auth"
	"net/http"

	"github.com/kataras/muxie"
)

// SetupAppServerV1 setup user-end routes.
func SetupAppServerV1(v1 muxie.SubMux) {
	{
		v1.Handle("/vauth/weapp", muxie.Methods().
			HandleFunc(http.MethodPost, app.User{}.AppWeappLogin))
		v1.Handle("/vauth/qapp", muxie.Methods().
			HandleFunc(http.MethodPost, app.User{}.AppQappLogin))

		v1.Handle("/vauth/signup", muxie.Methods().HandleFunc(http.MethodPost, app.User{}.SignUp))

		v1.Handle("/announces", muxie.Methods().HandleFunc(http.MethodGet, app.Announce{}.AppList))

		v1.Handle("/texts/:textID", muxie.Methods().HandleFunc(http.MethodGet, app.Text{}.AppView))

		v1.Handle("/topics", muxie.Methods().
			HandleFunc(http.MethodGet, app.Topic{}.List).
			HandleFunc(http.MethodPost, app.Topic{}.Create))
	}

	guard := muxie.Pre(auth.Guard("app", false))
	{
		v1.Handle("/status", muxie.Methods().Handle(http.MethodGet, guard.ForFunc(app.User{}.AppStatus)))

		v1.Handle("/tasks", muxie.Methods().
			Handle(http.MethodGet, guard.ForFunc(app.Task{}.AppList)))
		v1.Handle("/tasks/:taskID/bonus", muxie.Methods().
			Handle(http.MethodPost, guard.ForFunc(app.Task{}.AppTaskLogCreate)))
	}
}
