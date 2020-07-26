package route

import (
	"appsrv/admin"
	"appsrv/app"
	"net/http"

	"github.com/kataras/muxie"
)

// SetupAdminServerV1 setup admin routes.
func SetupAdminServerV1(v1 muxie.SubMux) {
	v1.Handle("/stat/summary", muxie.Methods().
		HandleFunc(http.MethodGet, app.Stat{}.Summary))

	v1.Handle(
		"/users",
		muxie.Methods().
			HandleFunc(http.MethodPost, admin.User{}.Create).
			HandleFunc(http.MethodGet, admin.User{}.List))

	v1.Handle("/tokens", muxie.Methods().
		HandleFunc(http.MethodPost, admin.Admin{}.Login))

	v1.Handle("/profile", muxie.Methods().
		HandleFunc(http.MethodGet, admin.Admin{}.Profile))

	v1.Handle("/admins", muxie.Methods().
		HandleFunc(http.MethodGet, admin.Admin{}.List).
		HandleFunc(http.MethodPost, admin.Admin{}.Create))

	v1.Handle("/admins/:id", muxie.Methods().
		HandleFunc(http.MethodDelete, admin.Admin{}.Delete))

	v1.Handle("/roles", muxie.Methods().
		HandleFunc(http.MethodGet, app.Role{}.List))

	v1.Handle("/admin-logs", muxie.Methods().HandleFunc(http.MethodGet, admin.AdminLog{}.List))

	v1.Handle("/password", muxie.Methods().HandleFunc(http.MethodPut, admin.Admin{}.UpdatePassword))

	v1.Handle("/texts", muxie.Methods().
		HandleFunc(http.MethodGet, app.Text{}.List).
		HandleFunc(http.MethodPost, app.Text{}.Create))
	v1.Handle("/texts/:textID", muxie.Methods().
		HandleFunc(http.MethodPut, app.Text{}.Update))

	v1.Handle("/announces", muxie.Methods().
		HandleFunc(http.MethodGet, app.Announce{}.List).
		HandleFunc(http.MethodPost, app.Announce{}.Create))
	v1.Handle("/announces/:announceID", muxie.Methods().
		HandleFunc(http.MethodDelete, app.Announce{}.Delete))

	v1.Handle("/tasks", muxie.Methods().
		HandleFunc(http.MethodGet, app.Task{}.List).
		HandleFunc(http.MethodPost, app.Task{}.Create))
	v1.Handle("/tasks/:taskID", muxie.Methods().
		HandleFunc(http.MethodDelete, app.Task{}.Delete))
}
