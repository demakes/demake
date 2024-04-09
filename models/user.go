package models

import (
	"github.com/gospel-sh/gospel/orm"
)

type User struct {
	orm.DBModel
	orm.JSONModel
	DisplayName string
	Source      string
	SourceID    string
	Superuser   bool
	EMail       string
}

type UserRole struct {
	orm.DBModel
	orm.JSONModel
	OrganizationID int64
	Organization   *Organization `db:"fk:OrganizationID"`
	UserID         int64
	User           *User `db:"fk:UserID"`
	Role           string
}
