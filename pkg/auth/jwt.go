package auth

import (
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/muxie"
)

func Token(scope string, uid uint) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.StandardClaims{
		Audience:  scope,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    fmt.Sprintf("%d", uid),
	})

	return t.SignedString([]byte(config.Server.Auth.Secret))
}

func Parse(token string) (*jwt.StandardClaims, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Server.Auth.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if t.Valid {
		if c, ok := t.Claims.(*jwt.StandardClaims); ok {
			return c, nil
		}

		return nil, errors.New("claims losed")
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
				t, err := jwt.Parse(token[7:], func(t *jwt.Token) (interface{}, error) {
					return []byte(config.Server.Auth.Secret), nil
				})

				if err != nil {
					if ve, ok := err.(*jwt.ValidationError); ok {
						if ve.Errors&jwt.ValidationErrorMalformed != 0 {
							if optional {
								next.ServeHTTP(w, r)
							} else {
								w.WriteHeader(http.StatusUnauthorized)
							}
						} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
							if optional {
								next.ServeHTTP(w, r)
							} else {
								w.WriteHeader(http.StatusUnauthorized)
							}
						} else {
							if optional {
								next.ServeHTTP(w, r)
							} else {
								w.WriteHeader(http.StatusUnauthorized)
							}
						}
					} else {
						if optional {
							next.ServeHTTP(w, r)
						} else {
							w.WriteHeader(http.StatusUnauthorized)
						}
					}
					return
				}

				if t.Valid {
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

	return db.DB.Model(ptr).Where("id = ?", t.Issuer).Select()
}
