package service

import (
	"appsrv/model"
	"appsrv/pkg/config"
	"appsrv/pkg/errors"
	"appsrv/schema"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/laeo/qapp"
	"github.com/medivhzhan/weapp/v2"
	"golang.org/x/crypto/bcrypt"
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

// LoginByWeapp 使用微信小程序进行登录
func (User) LoginByWeapp(db *pg.DB, conf config.WeappConfig, input schema.UserLoginByWeappInput) (*model.User, error) {
	res, err := weapp.Login(conf.AppID, conf.Secret, input.Code)
	if err != nil {
		return nil, errors.New(500, err.Error())
	}

	if err = res.GetResponseError(); err != nil {
		return nil, errors.New(500, err.Error())
	}

	var u model.User
	err = db.Model(&u).Where("open_id = ?", res.OpenID).First()
	if err != nil && pg.ErrNoRows != err {
		return nil, errors.New(500, err.Error())
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
			return nil, errors.New(500, err.Error())
		}
	}

	u.AvatarURL = input.UserInfo.AvatarURL
	_, err = db.Model(&u).Set("avatar_url = ?", u.AvatarURL).WherePK().Update()
	if err != nil {
		return nil, errors.New(500, err.Error())
	}

	return &u, nil
}

// LoginByQapp QQ小程序登录
func LoginByQapp(db *pg.DB, conf config.QappConfig, input schema.UserLoginByQappInput) (*model.User, error) {
	appid := conf.AppID
	secret := conf.Secret
	res, err := qapp.Login(appid, secret, input.Code)
	if err != nil {
		return nil, errors.New(500, err.Error())
	}

	if err := res.GetResponseError(); err != nil {
		return nil, errors.New(500, err.Error())
	}

	var u model.User
	err = db.Model(&u).Where("open_id = ?", res.OpenID).First()
	if err != nil && pg.ErrNoRows != err {
		return nil, errors.New(500, err.Error())
	}

	if err != nil && pg.ErrNoRows == err {
		// 创建用户
		u.AvatarURL = input.UserInfo.AvatarURL
		u.Name = input.UserInfo.NickName
		u.Gender = input.UserInfo.Gender
		u.OpenID = res.OpenID
		u.Remark = model.UserRemarkQQ
		err = db.Insert(&u)
		if err != nil {
			return nil, errors.New(500, err.Error())
		}
	}

	u.AvatarURL = input.UserInfo.AvatarURL
	_, err = db.Model(&u).Set("avatar_url = ?", u.AvatarURL).WherePK().Update()
	if err != nil {
		return nil, errors.New(500, err.Error())
	}

	return &u, nil
}

// SignByCredential 账号注册
func (User) SignByCredential(db *pg.DB, input schema.UserSignByCredentialInput) (*model.User, error) {
	var u model.User
	n, err := db.Model(&u).Where("name = ?", strings.ToLower(input.Username)).Count()
	if err != nil {
		return nil, errors.New(500, err.Error())
	}

	if n > 0 {
		return nil, errors.New(409, "该名称已被使用")
	}

	p, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(500, "网络通信错误请重试")
	}

	u.Name = input.Username
	u.Password = string(p)
	err = db.Insert(&u)
	if err != nil {
		return nil, errors.New(500, "网络通信错误请重试")
	}

	return &u, nil
}
