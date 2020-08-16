package auth

import (
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/out"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ContextKey 上下文键名
type ContextKey string

// Validated 用于从上下文提取JWT解析的数据
const Validated ContextKey = "validated"

// RoleClaims 添加了UserID和角色列表的JWT Claims.
type RoleClaims struct {
	jwt.StandardClaims

	UserID uint     `json:"uid"`
	Roles  []string `json:"ros"`
}

func Token(uid uint, roles []string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS384, RoleClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
		UserID: uid,
		Roles:  roles,
	})

	return t.SignedString([]byte(config.ProvideSingleton().Auth.Secret))
}

func GetUser(r *http.Request, ptr interface{}) error {
	rc := r.Context().Value(Validated)
	if rc == nil {
		return out.Err(401, "请登录")
	}

	dbc := db.ProvideSingleton()
	return dbc.Model(ptr).Where("id = ?", rc.(*RoleClaims).UserID).First()
}
