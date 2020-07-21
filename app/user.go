package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/errors"
	"appsrv/schema"
	"appsrv/service"
	"net/http"

	"github.com/kataras/muxie"
	"go.uber.org/zap"
)

type User struct{}

func (User) Create(w http.ResponseWriter, r *http.Request) {
	input := service.UserCreateInput{}
	err := muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, errors.ErrBodyBind)
		return
	}

	u, err := service.User{}.Create(db.DB, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	muxie.Dispatch(w, muxie.JSON, u)
}

func (User) List(w http.ResponseWriter, r *http.Request) {
	users, err := service.User{}.List(db.DB)
	if err != nil {
		bog.Error("User.List", zap.Error(err))
	}
	muxie.Dispatch(w, muxie.JSON, users)
}

func (User) AppWeappLogin(w http.ResponseWriter, r *http.Request) {
	var in schema.UserLoginByWeappInput
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	u, err := service.User{}.LoginByWeapp(db.DB, config.Server.Weapp, in)
	if err != nil {
		bog.Error("User.AppWeappLogin", zap.Error(err))
		w.WriteHeader(err.(errors.JSONError).Code)
		return
	}

	token, err := auth.Token("app", u.ID)
	if err != nil {
		bog.Error("User.AppWeappLogin", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	var out struct {
		Token string
		User  struct {
			ID        uint
			Name      string
			AvatarURL string
			Gender    int
			Created   string
			Coin      int
		}
	}
	out.Token = token
	out.User.ID = u.ID
	out.User.Name = u.Name
	out.User.AvatarURL = u.AvatarURL
	out.User.Gender = u.Gender
	out.User.Created = u.Created.Format("2006-01-02 15:04:05")

	err = muxie.Dispatch(w, muxie.JSON, &out)
	if err != nil {
		bog.Error("User.AppWeappLogin", zap.Error(err), zap.Uint("UID", u.ID))
		w.WriteHeader(500)
		return
	}
}

func (User) AppStatus(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var out = struct {
		ID        uint
		Name      string
		AvatarURL string
		Gender    int
		Coin      int
	}{
		ID:        u.ID,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Gender:    u.Gender,
		Coin:      u.Coin,
	}

	muxie.Dispatch(w, muxie.JSON, &out)
}

func (User) AppQappLogin(w http.ResponseWriter, r *http.Request) {
	var in schema.UserLoginByQappInput
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	u, err := service.LoginByQapp(db.DB, config.Server.Qapp, in)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err.(errors.JSONError))
		return
	}

	token, err := auth.Token("app", u.ID)
	if err != nil {
		bog.Error("User.AppQappLogin", zap.Error(err), zap.String("OpenID", u.OpenID))
		w.WriteHeader(500)
		return
	}

	var out struct {
		Token string
		User  struct {
			ID        uint
			Name      string
			AvatarURL string
			Gender    int
			Created   string
			Coin      int
		}
	}
	out.Token = token
	out.User.ID = u.ID
	out.User.Name = u.Name
	out.User.AvatarURL = u.AvatarURL
	out.User.Gender = u.Gender
	out.User.Created = u.Created.Format("2006-01-02 15:04:05")

	err = muxie.Dispatch(w, muxie.JSON, &out)
	if err != nil {
		bog.Error("User.AppQappLogin", zap.Error(err), zap.Uint("UID", u.ID))
		w.WriteHeader(500)
		return
	}
}

func (User) AppSignByCredential(w http.ResponseWriter, r *http.Request) {
	var input schema.UserSignByCredentialInput
	err := muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, errors.ErrBodyBind)
		return
	}

	u, err := service.User{}.SignByCredential(db.DB, input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err.(errors.JSONError))
		return
	}

	token, err := auth.Token("app", u.ID)
	if err != nil {
		bog.Error("User.AppQappLogin", zap.Error(err), zap.String("OpenID", u.OpenID))
		w.WriteHeader(500)
		return
	}

	var out struct {
		Token string
		User  struct {
			ID        uint
			Name      string
			AvatarURL string
			Gender    int
			Created   string
			Coin      int
		}
	}
	out.Token = token
	out.User.ID = u.ID
	out.User.Name = u.Name
	out.User.AvatarURL = u.AvatarURL
	out.User.Gender = u.Gender
	out.User.Created = u.Created.Format("2006-01-02 15:04:05")

	err = muxie.Dispatch(w, muxie.JSON, &out)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}
