package service

import (
	"appsrv/model"
	"appsrv/pkg/bog"
	"appsrv/pkg/errors"
	"appsrv/pkg/oss"
	"appsrv/schema"
	"bytes"
	"image/png"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/minio/minio-go/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/tsdtsdtsd/identicon"
	"go.uber.org/zap"
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
	input.Username = strings.TrimSpace(input.Username)

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
		return nil, errors.Err500
	}

	u.Name = input.Username
	u.Password = string(p)
	u.RoleKeys = []string{"v1"}
	u.LogoPath = time.Now().Format("20060102") + "/" + uuid.NewV4().String() + ".png"

	icon, err := identicon.New(u.Name, &identicon.Options{
		BackgroundColor: identicon.RGB(240, 240, 240),
		ImageSize:       240,
	})
	if err != nil {
		return nil, errors.Err500
	}
	b := bytes.NewBuffer(nil)
	err = png.Encode(b, icon)
	if err != nil {
		return nil, errors.Err500
	}

	_, err = oss.Min.PutObject(oss.Bucket, u.LogoPath, b, int64(b.Len()), minio.PutObjectOptions{
		ContentType: "image/png",
	})
	if err != nil {
		bog.Error("minio.PutObject", zap.Error(err))
		return nil, errors.Err500
	}

	err = db.Insert(&u)
	if err != nil {
		return nil, errors.Err500
	}

	return &u, nil
}

// FindWithCredential 根据凭据查找对应用户
func (User) FindWithCredential(db *pg.DB, input schema.UserSignUpInput) (*model.User, error) {
	var u model.User
	err := db.Model(&u).Where("name = ? OR name = ?", input.Username, strings.ToLower(input.Username)).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New(400, "该用户不存在")
		}
		return nil, errors.New(500, "服务器已爆炸")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New(400, "密码错误")
	}

	return &u, nil
}

// FindByName 根据给定名称查找对应用户
func (User) FindByName(db *pg.DB, name string) (*model.User, error) {
	var u model.User
	err := db.Model(&u).Where("name = ? OR name = ?", name, strings.ToLower(name)).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Err500
	}

	return &u, nil
}
