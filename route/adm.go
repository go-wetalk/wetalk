package route

import (
	"appsrv/admin"
	"appsrv/app"
	"appsrv/pkg/auth"
	"net/http"

	"github.com/kataras/muxie"
)

// SetupAdminServerV1 setup admin routes.
func SetupAdminServerV1(v1 muxie.SubMux) {
	guard := muxie.Pre(auth.Guard("admin", false))
	checkRoleRoot := muxie.Pre(app.Role{}.CheckRole("root"))
	checkForUser := muxie.Pre(app.Role{}.CheckRole("root", "user"))              // 管理用户
	checkForUserRO := muxie.Pre(app.Role{}.CheckRole("root", "user", "user:ro")) // 查询用户
	checkForText := muxie.Pre(app.Role{}.CheckRole("root", "text"))              // 文本管理
	checkForTextRO := muxie.Pre(app.Role{}.CheckRole("root", "text", "text:ro")) // 文本检索

	{
		v1.Handle("/stat/summary", muxie.Methods().
			Handle(http.MethodGet, guard.ForFunc(app.Stat{}.Summary)))

		v1.Handle(
			"/users",
			muxie.Methods().
				Handle(http.MethodPost, checkForUser.ForFunc(admin.User{}.Create)).
				Handle(http.MethodGet, checkForUserRO.ForFunc(admin.User{}.List)))

		v1.Handle("/tokens", muxie.Methods().
			HandleFunc(http.MethodPost, admin.Admin{}.Login))

		v1.Handle("/profile", muxie.Methods().
			Handle(http.MethodGet, guard.ForFunc(admin.Admin{}.Profile)))

		v1.Handle("/admins", muxie.Methods().
			Handle(http.MethodGet, checkRoleRoot.ForFunc(admin.Admin{}.List)).
			Handle(http.MethodPost, checkRoleRoot.ForFunc(admin.Admin{}.Create)))

		v1.Handle("/admins/:id", muxie.Methods().Handle(http.MethodDelete, checkRoleRoot.ForFunc(admin.Admin{}.Delete)))

		v1.Handle("/roles", muxie.Methods().
			Handle(http.MethodGet, checkRoleRoot.ForFunc(app.Role{}.List)))

		v1.Handle("/admin-logs", muxie.Methods().Handle(http.MethodGet, checkRoleRoot.ForFunc(admin.AdminLog{}.List)))

		v1.Handle("/password", muxie.Methods().Handle(http.MethodPut, guard.ForFunc(admin.Admin{}.UpdatePassword)))

		v1.Handle("/texts", muxie.Methods().
			Handle(http.MethodGet, checkForTextRO.ForFunc(app.Text{}.List)).
			Handle(http.MethodPost, checkForText.ForFunc(app.Text{}.Create)))
		v1.Handle("/texts/:textID", muxie.Methods().
			Handle(http.MethodPut, checkForText.ForFunc(app.Text{}.Update)))

		v1.Handle("/announces", muxie.Methods().
			Handle(http.MethodGet, guard.ForFunc(app.Announce{}.List)).
			Handle(http.MethodPost, guard.ForFunc(app.Announce{}.Create)))
		v1.Handle("/announces/:announceID", muxie.Methods().
			Handle(http.MethodDelete, guard.ForFunc(app.Announce{}.Delete)))

		v1.Handle("/tasks", muxie.Methods().
			Handle(http.MethodGet, guard.ForFunc(app.Task{}.List)).
			Handle(http.MethodPost, guard.ForFunc(app.Task{}.Create)))
		v1.Handle("/tasks/:taskID", muxie.Methods().
			Handle(http.MethodDelete, guard.ForFunc(app.Task{}.Delete)))
	}
}
