package service

import (
	"appsrv/model"
	"appsrv/pkg/config"
	"appsrv/pkg/out"
	"appsrv/schema"
	"bytes"
	"image/png"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/minio/minio-go/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/tsdtsdtsd/identicon"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// User 用户相关DB操作
type User struct {
	db   *pg.DB
	log  *zap.Logger
	mc   *minio.Client
	conf *config.ServerConfig
}

// CreateWithInput 账号注册
func (v *User) CreateWithInput(input schema.UserSignUpInput) (*model.User, error) {
	input.Username = strings.TrimSpace(input.Username)

	var u model.User
	n, err := v.db.Model(&u).Where("name = ?", strings.ToLower(input.Username)).Count()
	if err != nil {
		return nil, out.Err(500, err.Error())
	}

	if n > 0 {
		return nil, out.Err(409, "该名称已被使用")
	}

	p, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, out.Err500
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
		return nil, out.Err500
	}
	b := bytes.NewBuffer(nil)
	err = png.Encode(b, icon)
	if err != nil {
		return nil, out.Err500
	}

	_, err = v.mc.PutObject(v.conf.Oss.Bucket, u.LogoPath, b, int64(b.Len()), minio.PutObjectOptions{
		ContentType: "image/png",
	})
	if err != nil {
		v.log.Error("minio.PutObject", zap.Error(err))
		return nil, out.Err500
	}

	_, err = v.db.Model(&u).Insert()
	if err != nil {
		return nil, out.Err500
	}

	return &u, nil
}

// FindWithCredential 根据凭据查找对应用户
func (v *User) FindWithCredential(input schema.UserSignUpInput) (*model.User, error) {
	var u model.User
	err := v.db.Model(&u).Where("name = ? OR name = ?", input.Username, strings.ToLower(input.Username)).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, out.Err(400, "该用户不存在")
		}
		return nil, out.Err(500, "服务器已爆炸")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		return nil, out.Err(400, "密码错误")
	}

	return &u, nil
}

// FindByName 根据给定名称查找对应用户
func (v *User) FindByName(name string) (*model.User, error) {
	var u model.User
	err := v.db.Model(&u).Where("name = ? OR name = ?", name, strings.ToLower(name)).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, out.ErrNotFound
		}
		return nil, out.Err500
	}

	return &u, nil
}
