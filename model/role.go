package model

type Role struct {
	ID     uint
	Key    string
	Name   string
	Intro  string
	Admins []Admin `pg:"many2many:admin_roles,joinFK:role_id" json:"-"`
}
