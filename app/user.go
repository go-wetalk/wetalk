package app

import (
	"appsrv/model"
	"appsrv/pkg/auth"
	"appsrv/pkg/bog"
	"appsrv/pkg/config"
	"appsrv/pkg/db"
	"appsrv/pkg/errors"
	"appsrv/service"
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/kataras/muxie"
	"github.com/laeo/qapp"
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
		bog.Error("User.Create", zap.Error(err))
		muxie.Dispatch(w, muxie.JSON, errors.JSONError{
			Code: 500,
			Msg:  "保存失败",
		})
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
	var in service.UserLoginByWeappInput
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("User.AppWeappLogin", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	u, err := service.User{}.LoginByWeapp(db.DB, in)
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
	var in struct {
		Code          string `validate:"required" json:"code"`
		EncryptedData string `validate:"required" json:"encrypted_data"`
		IV            string `validate:"required" json:"iv"`
		RawData       string `validate:"required" json:"raw_data"`
		Signature     string `validate:"required" json:"signature"`
		UserInfo      struct {
			AvatarURL string `json:"avatarUrl"`
			NickName  string `json:"nickName"`
			Gender    int    `json:"gender"`
		} `validate:"required" json:"user_info"`
	}
	err := muxie.Bind(r, muxie.JSON, &in)
	if err != nil {
		bog.Error("User.AppQappLogin", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	appid := config.Server.Qapp.AppID
	secret := config.Server.Qapp.Secret
	res, err := qapp.Login(appid, secret, in.Code)
	if err != nil {
		bog.Error("User.AppQappLogin", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	if err := res.GetResponseError(); err != nil {
		bog.Error("User.AppQappLogin", zap.Error(err), zap.String("NickName", in.UserInfo.NickName), zap.String("appid", appid))
		w.WriteHeader(500)
		return
	}

	var u model.User
	err = db.DB.Model(&u).Where("open_id = ?", res.OpenID).First()
	if err != nil && pg.ErrNoRows != err {
		bog.Error("User.AppQappLogin", zap.Error(err), zap.String("OpenID", res.OpenID))
		w.WriteHeader(500)
		return
	}

	if err != nil && pg.ErrNoRows == err {
		// 创建用户
		u.AvatarURL = in.UserInfo.AvatarURL
		u.Name = in.UserInfo.NickName
		u.Gender = in.UserInfo.Gender
		u.OpenID = res.OpenID
		u.Remark = model.UserRemarkQQ
		err = db.DB.Insert(&u)
		if err != nil {
			bog.Error("User.AppQappLogin", zap.Error(err), zap.String("OpenID", res.OpenID))
			w.WriteHeader(500)
			return
		}
	}

	u.AvatarURL = in.UserInfo.AvatarURL
	_, err = db.DB.Model(&u).Set("avatar_url = ?", u.AvatarURL).WherePK().Update()
	if err != nil {
		bog.Error("User.AppQappLogin", zap.Error(err), zap.String("OpenID", res.OpenID))
		w.WriteHeader(500)
		return
	}

	token, err := auth.Token("app", u.ID)
	if err != nil {
		bog.Error("User.AppQappLogin", zap.Error(err), zap.String("OpenID", res.OpenID))
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
