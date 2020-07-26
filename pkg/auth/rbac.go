package auth

import (
	"appsrv/model"
	"appsrv/pkg/config"
	"appsrv/pkg/errors"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v9"
	"github.com/kataras/muxie"
	"github.com/storyicon/grbac"
)

// RoleGuard 角色检查中间件
func RoleGuard(db *pg.DB) muxie.Wrapper {
	rbac, err := grbac.New(grbac.WithLoader(RoleRulesLoader(db), time.Minute))
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assigned := []string{}
			if token := r.Header.Get("Authorization"); strings.HasPrefix(token, "Bearer ") {
				t, err := jwt.ParseWithClaims(token[7:], &RoleClaims{}, func(t *jwt.Token) (interface{}, error) {
					return []byte(config.Server.Auth.Secret), nil
				})

				if err == nil && t.Valid {
					if rc := t.Claims.(*RoleClaims); rc != nil {
						r = r.WithContext(context.WithValue(r.Context(), Validated, rc))
						assigned = append(assigned, rc.Roles...)
					}
				}
			}

			state, _ := rbac.IsRequestGranted(r, assigned)
			if state.IsGranted() || strings.Contains(strings.Join(assigned, ","), "root") {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(403)
				muxie.Dispatch(w, muxie.JSON, errors.New(403, "您无权访问该对象"))
			}
		})
	}
}

// RoleRulesLoader 角色权限配置加载器
func RoleRulesLoader(db *pg.DB) func() (rules grbac.Rules, err error) {
	return func() (rules grbac.Rules, err error) {
		rs := []model.Rule{}
		err = db.Model(&rs).Relation("AuthorizedRoles.key").Relation("ForbiddenRoles.key").Order("rule.priority ASC").Select()
		if err != nil {
			return nil, err
		}

		for _, r := range rs {
			rule := grbac.Rule{}
			rule.ID = int(r.ID)
			rule.Host = r.Host
			rule.Path = r.Path
			rule.Method = r.Method
			rule.AllowAnyone = r.AllowAnyone
			rule.AuthorizedRoles = r.Authorized
			rule.ForbiddenRoles = r.Forbidden
			rules = append(rules, &rule)
		}

		return
	}
}
