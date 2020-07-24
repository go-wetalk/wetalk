package service

import (
	"appsrv/model"
	"appsrv/pkg/errors"
	"appsrv/schema"
	"strings"

	"github.com/go-pg/pg/v9"
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

// CreateWithInput 账号注册
func (User) CreateWithInput(db *pg.DB, input schema.UserSignUpInput) (*model.User, error) {
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
