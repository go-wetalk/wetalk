package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/oss"
	"appsrv/pkg/out"
	"appsrv/schema"
	"appsrv/service"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/hcaptcha"
	"github.com/kataras/muxie"
	"github.com/minio/minio-go/v6"
	"github.com/xeonx/timeago"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct{}

func (User) AppStatus(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
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
		UnreadNotify: u.UnreadNotify(),
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(us))
}

func (User) SignUp(w http.ResponseWriter, r *http.Request) {
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

	if config.Server.HCaptcha.Enabled {
		hcc := hcaptcha.New(config.Server.HCaptcha.Secret)
		if resp := hcc.VerifyToken(input.Captcha); !resp.Success {
			muxie.Dispatch(w, muxie.JSON, out.Err(400, resp.ChallengeTS))
			return
		}
	}

	u, err := service.User{}.CreateWithInput(db.DB, input)
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
	raw.User.UnreadNotify = u.UnreadNotify()
	muxie.Dispatch(w, muxie.JSON, out.Data(raw))
}

func (User) Login(w http.ResponseWriter, r *http.Request) {
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

	u, err := service.User{}.FindWithCredential(db.DB, input)
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
	raw.User.UnreadNotify = u.UnreadNotify()
	muxie.Dispatch(w, muxie.JSON, out.Data(raw))
}

func (User) ViewUserDetail(w http.ResponseWriter, r *http.Request) {
	name := muxie.GetParam(w, "name")
	u, err := service.User{}.FindByName(db.DB, strings.TrimSpace(name))
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

func (User) ViewProfile(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := auth.GetUser(r, &u)
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err401)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(u))
}

func (User) UpdateLogo(w http.ResponseWriter, r *http.Request) {
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

	_, err = oss.Min.PutObject(oss.Bucket, objectPath, f, h.Size, minio.PutObjectOptions{
		ContentType: h.Header.Get("Content-Type"),
	})
	if err != nil {
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	_, err = db.DB.Model(&u).WherePK().Set("logo_path = ?", objectPath).Update()
	if err != nil {
		bog.Error(err.Error(), zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(nil))
}

func (User) UpdateAddress(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.DB.Model(&u).WherePK().
		Set("country = ?, province = ?, city = ?, street = ?", input.Country, input.Province, input.City, input.Street).
		Update()
	if err != nil {
		bog.Error(err.Error(), zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(nil))
}

func (User) UpdateSocial(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.DB.Model(&u).WherePK().
		Set("sign = ?", input.Sign).
		Update()
	if err != nil {
		bog.Error(err.Error(), zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(nil))
}

func (User) UpdatePassword(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.DB.Model(&u).WherePK().
		Set("password = ?", string(b)).
		Update()
	if err != nil {
		bog.Error(err.Error(), zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, out.Err500)
		return
	}

	muxie.Dispatch(w, muxie.JSON, out.Data(nil))
}
