package auth

import (
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/muxie"
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

	return t.SignedString([]byte(config.Server.Auth.Secret))
}

func Parse(token string) (*RoleClaims, error) {
	t, err := jwt.ParseWithClaims(token, &RoleClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Server.Auth.Secret), nil
	})

	if err == nil && t.Valid {
		return t.Claims.(*RoleClaims), nil
	}

	return nil, err
}

func Guard(scope string, optional bool) muxie.Wrapper {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if len(token) < 7 {
				if optional {
					next.ServeHTTP(w, r)
					return
				}

				w.WriteHeader(http.StatusUnauthorized)
			} else {
				t, err := jwt.ParseWithClaims(token[7:], &RoleClaims{}, func(t *jwt.Token) (interface{}, error) {
					return []byte(config.Server.Auth.Secret), nil
				})

				if err == nil && t.Valid {
					r = r.WithContext(context.WithValue(r.Context(), Validated, t.Claims.(*RoleClaims)))
					next.ServeHTTP(w, r)
				} else {
					if optional {
						next.ServeHTTP(w, r)
					} else {
						w.WriteHeader(http.StatusUnauthorized)
					}
				}
			}
		})
	}
}

func GetUser(r *http.Request, ptr interface{}) error {
	token := r.Header.Get("Authorization")
	if len(token) < 7 {
		return errors.New("请登录")
	}

	t, err := Parse(token[7:])
	if err != nil {
		return err
	}

	return db.DB.Model(ptr).Where("id = ?", t.UserID).Select()
}
