package models

import (
	"github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
)

type Site struct {
	orm.DBModel `db:"table:project"`
	orm.JSONModel
	HeadID      *int64 `db:"head_id"`
	Name        string
	Hostname    string
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

func (c *Site) ByHostname(hostname string) error {
	return orm.LoadOne(c, map[string]any{"hostname": hostname})
}

type SitePlugin interface {
}

type BlogPost struct {
	Title TranslatedString `json:"title"`
}

type BlogPlugin struct {
	ArticlesPerPage int         `json:"articlesPerPage"`
	Posts           []*BlogPost `json:"posts"`
}

type SiteMeta struct {
	Title  TranslatedString `json:"title"`
	Domain string           `json:"domain"`
}

type TranslatedString struct {
	Translations map[string]string `json:"translations"`
	ID           string            `json:"id"`
}

type SiteGraph struct {
	Plugins []SitePlugin       `json:"plugins" graph:"include"`
	Meta    SiteMeta           `json:"meta"`
	DOM     gospel.HTMLElement `json:"dom"`
}

func init() {
	MustRegister[SiteGraph]("siteGraph")
	MustRegister[TranslatedString]("translatedString")
	MustRegister[SiteMeta]("siteMeta")
	MustRegister[BlogPlugin]("blogPlugin")
	MustRegister[BlogPost]("blogPost")
}
