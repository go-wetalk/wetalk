package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/config"
	"appsrv/pkg/out"
	"appsrv/schema"
	"appsrv/service"
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/kataras/hcaptcha"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"github.com/xeonx/timeago"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig

	userService *service.User
}

func (v *User) RegisterRoute(m muxie.SubMux) {
	m.Handle("/users", muxie.Methods().HandleFunc(http.MethodPost, v.SignUp))
	m.Handle("/tokens", muxie.Methods().HandleFunc(http.MethodPost, v.Login))
	m.Handle("/status", muxie.Methods().HandleFunc(http.MethodGet, v.AppStatus))
	m.Handle("/users/:name", muxie.Methods().HandleFunc(http.MethodGet, v.ViewUserDetail))
	m.Handle("/profile", muxie.Methods().HandleFunc(http.MethodGet, v.ViewProfile))
	m.Handle("/profile/logo", muxie.Methods().HandleFunc(http.MethodPut, v.UpdateLogo))
	m.Handle("/profile/address", muxie.Methods().HandleFunc(http.MethodPut, v.UpdateAddress))
	m.Handle("/profile/social", muxie.Methods().HandleFunc(http.MethodPut, v.UpdateSocial))
	m.Handle("/profile/password", muxie.Methods().HandleFunc(http.MethodPut, v.UpdatePassword))
}

func (v *User) AppStatus(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	us := schema.UserStatus{
		ID:           u.ID,
		Name:         u.Name,
		Logo:         u.Logo,
		Gender:       u.Gender,
		Coin:         u.Coin,
		Created:      u.Created.Format("2006-01-02 15:04:05"),
		RoleList:     u.RoleKeys,
		UnreadNotify: u.UnreadNotify(v.db),
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(us))
}

func (v *User) SignUp(w http.ResponseWriter, r *http.Request) {
	var input schema.UserSignUpInput
	err := muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	if err = input.Validate(); err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err(400, err.Error()))
		return
	}

	if v.conf.HCaptcha.Enabled {
		hcc := hcaptcha.New(v.conf.HCaptcha.Secret)
		if resp := hcc.VerifyToken(input.Captcha); !resp.Success {
			muxie.Dispatch(w, muxie.JSON, out.Err(400, resp.ChallengeTS))
			return
		}
	}

	u, err := v.userService.CreateWithInput(input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	token, err := auth.Token(u.ID, u.RoleKeys)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	var raw struct {
		Token string
		User  schema.UserStatus
	}
	raw.Token = token
	raw.User.ID = u.ID
	raw.User.Name = u.Name
	raw.User.Logo = u.Logo
	raw.User.Gender = u.Gender
	raw.User.Created = u.Created.Format("2006-01-02 15:04:05")
	raw.User.RoleList = u.RoleKeys
	raw.User.UnreadNotify = u.UnreadNotify(v.db)
	muxie.Dispatch(w, muxie.JSON, out.Data(raw))
}

func (v *User) Login(w http.ResponseWriter, r *http.Request) {
	var input schema.UserSignUpInput
	err := muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	if err = input.Validate(); err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err(400, err.Error()))
		return
	}

	u, err := v.userService.FindWithCredential(input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	token, err := auth.Token(u.ID, u.RoleKeys)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	var raw struct {
		Token string
		User  schema.UserStatus
	}
	raw.Token = token
	raw.User.ID = u.ID
	raw.User.Name = u.Name
	raw.User.Logo = u.Logo
	raw.User.Gender = u.Gender
	raw.User.Created = u.Created.Format("2006-01-02 15:04:05")
	raw.User.RoleList = u.RoleKeys
	raw.User.UnreadNotify = u.UnreadNotify(v.db)
	muxie.Dispatch(w, muxie.JSON, out.Data(raw))
}

func (v *User) ViewUserDetail(w http.ResponseWriter, r *http.Request) {
	name := muxie.GetParam(w, "name")
	u, err := v.userService.FindByName(strings.TrimSpace(name))
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, err)
		return
	}

	raw := schema.UserDetail{}
	raw.ID = u.ID
	raw.Name = u.Name
	raw.Logo = u.LogoLink()
	raw.Created = timeago.Chinese.Format(u.Created)

	muxie.Dispatch(w, muxie.JSON, out.Data(raw))
}

func (v *User) ViewProfile(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err401)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(u))
}

func (v *User) UpdateLogo(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err401)
		return
	}

	f, h, err := r.FormFile("logo")
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	objectPath := time.Now().Format("20060102") + "/" + h.Filename

	_, err = v.mc.PutObject(v.conf.Oss.Bucket, objectPath, f, h.Size, minio.PutObjectOptions{
		ContentType: h.Header.Get("Content-Type"),
	})
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	_, err = v.db.Model(&u).WherePK().Set("logo_path = ?", objectPath).Update()
	if err != nil {
		v.log.Error(err.Error(), zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(nil))
}

func (v *User) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err401)
		return
	}

	input := schema.UserAddress{}
	err = muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	_, err = v.db.Model(&u).WherePK().
		Set("country = ?, province = ?, city = ?, street = ?", input.Country, input.Province, input.City, input.Street).
		Update()
	if err != nil {
		v.log.Error(err.Error(), zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(nil))
}

func (v *User) UpdateSocial(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err401)
		return
	}

	input := schema.UserSocial{}
	err = muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	_, err = v.db.Model(&u).WherePK().
		Set("sign = ?", input.Sign).
		Update()
	if err != nil {
		v.log.Error(err.Error(), zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(nil))
}

func (v *User) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err401)
		return
	}

	input := schema.UserPassword{}
	err = muxie.Bind(r, muxie.JSON, &input)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.ErrBodyBind)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Original))
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err(400, "原始密码错误"))
		return
	}

	b, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err(500, "密码修改失败"))
		return
	}

	_, err = v.db.Model(&u).WherePK().
		Set("password = ?", string(b)).
		Update()
	if err != nil {
		v.log.Error(err.Error(), zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(nil))
}
