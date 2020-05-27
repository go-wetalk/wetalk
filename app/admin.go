package app

import (
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/db"
	"appsrv/pkg/errors"
	"net/http"
	"time"

	"github.com/kataras/muxie"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       uint
	Name     string
	Password string `json:"-"`
	Created  time.Time
	Updated  time.Time
	Deleted  *time.Time `pg:",soft_delete"`
	Roles    []Role     `pg:"many2many:admin_roles,joinFK:admin_id"`
}

func (a *Admin) RoleList() (roles []Role) {
	q := "select a.* from roles a join admin_roles b on b.admin_id = ? and b.role_id = a.id"
	_, err := db.DB.Query(&roles, q, a.ID)
	if err != nil {
		bog.Error("Admin.RoleList", zap.Error(err))
	}
	return
}

func (Admin) Login(w http.ResponseWriter, r *http.Request) {
	var cdt struct {
		Username string
		Password string
	}

	err := muxie.Bind(r, muxie.JSON, &cdt)
	if err != nil {
		bog.Error("Admin.Login", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		muxie.Dispatch(w, muxie.JSON, errors.ErrBodyBind)
		return
	}

	var a Admin
	err = db.DB.Model(&a).Where("name = ?", cdt.Username).Select()
	if err != nil {
		bog.Error("Admin.Login", zap.Error(err))
		a.Name = cdt.Username
		AdminLog{}.LogEvent(r, &a, "login", "账号 "+cdt.Username+" 不存在")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(cdt.Password))
	if err != nil {
		AdminLog{}.LogEvent(r, &a, "login", "密码错误")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	AdminLog{}.LogEvent(r, &a, "login", "成功")

	token, err := auth.Token("admin", a.ID)
	if err != nil {
		bog.Error("Admin.Login", zap.Error(err))
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	var bind struct {
		Token string
	}
	bind.Token = token
	err = muxie.Dispatch(w, muxie.JSON, &bind)
	if err != nil {
		bog.Error("Admin.Login", zap.Error(err))
	}
}

func (Admin) Profile(w http.ResponseWriter, r *http.Request) {
	a := Admin{}
	err := auth.GetUser(r, &a)
	if err != nil {
		bog.Error("Admin.Profile", zap.Error(err))
		w.WriteHeader(401)
		muxie.Dispatch(w, muxie.JSON, errors.JSONError{
			Code: 401,
			Msg:  err.Error(),
		})
		return
	}

	a.Roles = a.RoleList()

	err = muxie.Dispatch(w, muxie.JSON, &a)
	if err != nil {
		bog.Error("Admin.Profile", zap.Error(err))
		w.WriteHeader(500)
	}
}

func (Admin) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Admin
		Password string
	}

	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Admin.Create", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	var a = in.Admin

	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	a.Password = string(hash)

	err = db.DB.Insert(&a)
	if err != nil {
		bog.Error("Admin.Create", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	for _, role := range a.Roles {
		_, err = db.DB.ExecOne("insert into admin_roles(admin_id, role_id) values(?, ?)", a.ID, role.ID)
		if err != nil {
			bog.Error("Admin.Create", zap.Error(err))
			w.WriteHeader(500)
			return
		}
	}

	muxie.Dispatch(w, muxie.JSON, &a)
}

func (Admin) List(w http.ResponseWriter, r *http.Request) {
	var as = []Admin{}
	err := db.DB.Model(&as).Select()
	if err != nil {
		bog.Error("Admin.List", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	for i := range as {
		as[i].Roles = as[i].RoleList()
	}

	muxie.Dispatch(w, muxie.JSON, &as)
}

func (Admin) Delete(w http.ResponseWriter, r *http.Request) {
	id := muxie.GetParam(w, "id")
	var a Admin
	err := db.DB.Model(&a).Where("id = ?", id).First()
	if err != nil {
		bog.Error("Admin.Delete", zap.Error(err))
		w.WriteHeader(404)
		return
	}

	err = db.DB.Delete(&a)
	if err != nil {
		bog.Error("Admin.Delete", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (Admin) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var in struct {
		OldPwd  string
		NewPwd  string
		Confirm string
	}

	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("Admin.UpdatePassword", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	if in.NewPwd != in.Confirm {
		muxie.Dispatch(w, muxie.JSON, errors.JSONError{
			Code: 400,
			Msg:  "新密码输入不一致",
		})
		return
	}

	var a Admin
	err = auth.GetUser(r, &a)
	if err != nil {
		bog.Error("Admin.UpdatePassword", zap.Error(err))
		w.WriteHeader(401)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(in.OldPwd))
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, errors.JSONError{
			Code: 400,
			Msg:  "密码错误",
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(in.NewPwd), bcrypt.DefaultCost)
	_, err = db.DB.Model(&a).Set("password = ?", string(hash)).WherePK().Update()
	if err != nil {
		bog.Error("Admin.UpdatePassword", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	AdminLog{}.LogEvent(r, &a, "password", "修改密码")

	muxie.Dispatch(w, muxie.JSON, errors.JSONError{
		Code: 200,
		Msg:  "修改成功",
	})
}
