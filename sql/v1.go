package sql

import (
	"time"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	Setup(1, func(db *pg.DB, l *zap.Logger) error {
		createTable(
			&admin{},
			&adminLog{},
			&adminRole{},
			&role{},
		)

		i, _ := db.Model(&admin{}).Count()
		if i == 0 {
			l.Info("init table admin", zap.String("username", "admin"), zap.String("password", "admina"))
			hash, _ := bcrypt.GenerateFromPassword([]byte("admina"), bcrypt.DefaultCost)
			db.Insert(&admin{
				Name:     "admin",
				Password: string(hash),
			})
		}

		i, _ = db.Model(&role{}).Count()
		if i == 0 {
			l.Info("init table role", zap.String("key", "root"))
			db.Insert(&role{
				Key:  "root",
				Name: "系统管理员",
			})
		}

		i, _ = db.Model(&adminRole{}).Count()
		if i == 0 {
			l.Info("init role binding", zap.String("admin", "id:1"), zap.String("role", "key:root"))
			db.Insert(&adminRole{
				AdminID: 1,
				RoleID:  1,
			})
		}

		return nil
	})
}

type admin struct {
	ID       uint
	Name     string
	Password string    `json:"-"`
	Created  time.Time `pg:",default:now()"`
	Updated  time.Time `pg:",default:now()"`
	Deleted  *time.Time
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
	ID     uint
	Key    string `pg:",unique"`
	Name   string
	Intro  string
	Admins []admin `pg:"many2many:admin_roles,joinFK:role_id" json:"-"`
}

type adminRole struct {
	AdminID uint
	RoleID  uint
}
