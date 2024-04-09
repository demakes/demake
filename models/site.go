package models

import (
	"github.com/gospel-sh/gospel/orm"
)

type Site struct {
	orm.DBModel `db:"table:project"`
	orm.JSONModel
	Name        string
	Description string
}

func (c *Site) Save() error {
	return orm.Save(c)
}

func (c *Site) ByExtID(id []byte) error {
	return orm.LoadOne(c, map[string]any{"ext_id": id})
}

func (c *Site) ByID(id int64) error {
	return orm.LoadOne(c, map[string]any{"id": id})
}
