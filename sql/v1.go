package sql

import (
	"appsrv/pkg/db"
	"time"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	Setup(1, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			db,
			&admin{},
			&adminLog{},
			&role{},
			&rule{},
		)

		db.Model(
			[]*role{
				&role{Key: "v0", Name: "V0·超管"},
				&role{Key: "v1", Name: "V1·用户"},
			}).Insert()

		l.Info("init table admin", zap.String("username", "admin"), zap.String("password", "admina"))
		hash, _ := bcrypt.GenerateFromPassword([]byte("admina"), bcrypt.DefaultCost)
		db.Model(&admin{
			Name:     "admin",
			Password: string(hash),
			RoleKeys: []string{"v0"},
		}).Insert()

		db.Model(
			// 游客级规则拥有最低优先级
			&rule{
				Host:        "*",
				Path:        "/*", // 接口前部必定由版本号起头，所以需要斜杠来分割节点，否则匹配的规则会被组合成类似 /v[0-9]* 的形式导致匹配失败
				Method:      "{HEAD,GET}",
				AllowAnyone: true,
			},
		).Insert()

		return nil
	})
}

type admin struct {
	ID       uint
	Name     string   `pg:",unique"`
	Password string   `json:"-"`
	RoleKeys []string `pg:",array"`

	db.TimeUpdate
}

type adminLog struct {
	ID        uint
	AdminID   uint
	AdminName string
	Event     string
	Intro     string
	IP        string
	UA        string
	Ref       string
	Created   time.Time `pg:",default:now()"`
}

type role struct {
	ID    uint
	Key   string `pg:",unique"`
	Name  string
	Intro string

	db.TimeUpdate
}

type rule struct {
	ID          uint
	Host        string `pg:",unique:action,default:'*'"`
	Path        string `pg:",unique:action,default:'*'"`
	Method      string `pg:",unique:action,default:'*'"`
	AllowAnyone bool
	Authorized  []string `pg:",array"`
	Forbidden   []string `pg:",array"`

	db.TimeUpdate
}
