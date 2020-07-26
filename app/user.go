package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/db"
	"appsrv/pkg/errors"
	"appsrv/schema"
	"appsrv/service"
	"net/http"

	"github.com/kataras/muxie"
)

type User struct{}

func (User) AppStatus(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var out = struct {
		ID     uint
		Name   string
		Logo   string
		Gender int
		Coin   int
	}{
		ID:     u.ID,
		Name:   u.Name,
		Logo:   u.Logo,
		Gender: u.Gender,
		Coin:   u.Coin,
	}

	muxie.Dispatch(w, muxie.JSON, &out)
}

func (User) SignUp(w http.ResponseWriter, r *http.Request) {
	var input schema.UserSignUpInput
	err := muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		w.WriteHeader(429)
		muxie.Dispatch(w, muxie.JSON, errors.ErrBodyBind)
		return
	}

	if err = input.Validate(); err != nil {
		w.WriteHeader(400)
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	u, err := service.User{}.CreateWithInput(db.DB, input)
	if err != nil {
		w.WriteHeader(500)
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	token, err := auth.Token("app", u.ID, []string{})
	if err != nil {
		w.WriteHeader(500)
		muxie.Dispatch(w, muxie.JSON, errors.New(500, "网络链接波动请重试"))
		return
	}

	var out struct {
		Token string
		User  struct {
			ID      uint
			Name    string
			Logo    string
			Gender  int
			Created string
			Coin    int
		}
	}
	out.Token = token
	out.User.ID = u.ID
	out.User.Name = u.Name
	out.User.Logo = u.Logo
	out.User.Gender = u.Gender
	out.User.Created = u.Created.Format("2006-01-02 15:04:05")
	muxie.Dispatch(w, muxie.JSON, out)
}

func (User) Login(w http.ResponseWriter, r *http.Request) {
	var input schema.UserSignUpInput
	err := muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		w.WriteHeader(429)
		muxie.Dispatch(w, muxie.JSON, errors.ErrBodyBind)
		return
	}

	if err = input.Validate(); err != nil {
		w.WriteHeader(400)
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	u, err := service.User{}.FindWithCredential(db.DB, input)
	if err != nil {
		w.WriteHeader(err.(errors.JSONError).Code)
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	token, err := auth.Token("app", u.ID, []string{})
	if err != nil {
		w.WriteHeader(500)
		muxie.Dispatch(w, muxie.JSON, errors.New(500, "网络链接波动请重试"))
		return
	}

	var out struct {
		Token string
		User  struct {
			ID      uint
			Name    string
			Logo    string
			Gender  int
			Created string
			Coin    int
		}
	}
	out.Token = token
	out.User.ID = u.ID
	out.User.Name = u.Name
	out.User.Logo = u.Logo
	out.User.Gender = u.Gender
	out.User.Created = u.Created.Format("2006-01-02 15:04:05")
	muxie.Dispatch(w, muxie.JSON, out)
}
