package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/errors"
	"appsrv/schema"
	"appsrv/service"
	"net/http"
	"strings"

	"github.com/kataras/hcaptcha"
	"github.com/kataras/muxie"
	"github.com/xeonx/timeago"
)

type User struct{}

func (User) AppStatus(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	var out = schema.UserStatus{
		ID:           u.ID,
		Name:         u.Name,
		Logo:         u.Logo,
		Gender:       u.Gender,
		Coin:         u.Coin,
		Created:      u.Created.Format("2006-01-02 15:04:05"),
		RoleList:     u.RoleKeys,
		UnreadNotify: u.UnreadNotify(),
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

	if config.Server.HCaptcha.Enabled {
		hcc := hcaptcha.New(config.Server.HCaptcha.Secret)
		if resp := hcc.VerifyToken(input.Captcha); !resp.Success {
			w.WriteHeader(400)
			muxie.Dispatch(w, muxie.JSON, errors.New(400, resp.ChallengeTS))
			return
		}
	}

	u, err := service.User{}.CreateWithInput(db.DB, input)
	if err != nil {
		w.WriteHeader(500)
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	token, err := auth.Token(u.ID, u.RoleKeys)
	if err != nil {
		w.WriteHeader(500)
		muxie.Dispatch(w, muxie.JSON, errors.New(500, "网络链接波动请重试"))
		return
	}

	var out struct {
		Token string
		User  schema.UserStatus
	}
	out.Token = token
	out.User.ID = u.ID
	out.User.Name = u.Name
	out.User.Logo = u.Logo
	out.User.Gender = u.Gender
	out.User.Created = u.Created.Format("2006-01-02 15:04:05")
	out.User.RoleList = u.RoleKeys
	out.User.UnreadNotify = u.UnreadNotify()
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

	token, err := auth.Token(u.ID, u.RoleKeys)
	if err != nil {
		w.WriteHeader(500)
		muxie.Dispatch(w, muxie.JSON, errors.New(500, "网络链接波动请重试"))
		return
	}

	var out struct {
		Token string
		User  schema.UserStatus
	}
	out.Token = token
	out.User.ID = u.ID
	out.User.Name = u.Name
	out.User.Logo = u.Logo
	out.User.Gender = u.Gender
	out.User.Created = u.Created.Format("2006-01-02 15:04:05")
	out.User.RoleList = u.RoleKeys
	out.User.UnreadNotify = u.UnreadNotify()
	muxie.Dispatch(w, muxie.JSON, out)
}

func (User) ViewUserDetail(w http.ResponseWriter, r *http.Request) {
	name := muxie.GetParam(w, "name")
	u, err := service.User{}.FindByName(db.DB, strings.TrimSpace(name))
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	out := schema.UserDetail{}
	out.ID = u.ID
	out.Name = u.Name
	out.Logo = u.LogoLink()
	out.Created = timeago.Chinese.Format(u.Created)

	muxie.Dispatch(w, muxie.JSON, out)
}
