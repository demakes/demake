package models

import (
	"github.com/gospel-sh/gospel/orm"
)

type Organization struct {
	orm.DBModel
	orm.JSONModel
	Name        string
	Source      string
	SourceID    []byte
	Description string
}
