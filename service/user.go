package service

import (
	"appsrv/model"
	"appsrv/pkg/config"
	"appsrv/pkg/errors"

	"github.com/go-pg/pg/v9"
	"github.com/medivhzhan/weapp/v2"
)

type User struct{}

type UserCreateInput model.User

func (User) Create(db *pg.DB, u UserCreateInput) (model.User, error) {
	err := db.Insert(&u)
	return model.User(u), err
}

func (User) List(db *pg.DB) ([]model.User, error) {
	var users = []model.User{}
	err := db.Model(&users).Order("id ASC").Select()
	return users, err
}

// UserLoginByWeappInput 微信小程序登录参数
type UserLoginByWeappInput struct {
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

// LoginByWeapp 使用微信小程序进行登录
func (User) LoginByWeapp(db *pg.DB, input UserLoginByWeappInput) (*model.User, error) {
	res, err := weapp.Login(config.Server.Weapp.AppID, config.Server.Weapp.Secret, input.Code)
	if err != nil {
		return nil, errors.JSONError{Code: 500, Msg: err.Error()}
	}

	if err = res.GetResponseError(); err != nil {
		return nil, errors.JSONError{Code: 500, Msg: err.Error()}
	}

	var u model.User
	err = db.Model(&u).Where("open_id = ?", res.OpenID).First()
	if err != nil && pg.ErrNoRows != err {
		return nil, errors.JSONError{Code: 500, Msg: err.Error()}
	}

	if err != nil && pg.ErrNoRows == err {
		// 创建用户
		u.AvatarURL = input.UserInfo.AvatarURL
		u.Name = input.UserInfo.NickName
		u.Gender = input.UserInfo.Gender
		u.OpenID = res.OpenID
		u.Remark = model.UserRemarkWechat
		err = db.Insert(&u)
		if err != nil {
			return nil, errors.JSONError{Code: 500, Msg: err.Error()}
		}
	}

	u.AvatarURL = input.UserInfo.AvatarURL
	_, err = db.Model(&u).Set("avatar_url = ?", u.AvatarURL).WherePK().Update()
	if err != nil {
		return nil, errors.JSONError{Code: 500, Msg: err.Error()}
	}

	return &u, nil
}
