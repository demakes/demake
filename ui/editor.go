package ui

import (
	"encoding/hex"
	"fmt"
	. "github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
	"github.com/demakes/demake/models"
)

func EditSite(c Context, siteID string) Element {

	id, err := hex.DecodeString(siteID)
	db := UseDB(c)
	dbf := func() orm.DB { return db }

	if err != nil {
		return Div("invalid ID")
	}

	if UseUser(c) == nil {
		UseRouter(c).RedirectTo("/login")
		return nil
	}

	site := orm.Init(&models.Site{}, dbf)

	if err := site.ByExtID(id); err != nil {
		return Div("cannot find site")
	}

	siteGraph, err := GetGraph(site, dbf)

	if err != nil {
		return Div(Fmt("Cannot load site: %v", err))
	}

	return SiteEditor(c, site, siteGraph, dbf)
}

func SiteEditor(c Context, site *models.Site, siteGraph *models.SiteGraph, dbf func() orm.DB) Element {

	form := MakeFormData(c, "editor", POST)
	source := form.Var("source", siteGraph.DOM.RenderCode())
	router := UseRouter(c)
	error := Var(c, "")

	onSubmit := func() {

		parser := &Parser{
			Source: source.Get(),
		}

		element, err := parser.ParseHTMLElement()

		if err != nil {
			error.Set(Fmt("cannot parse: %v", err))
			return
		}

		if element == nil {
			error.Set("not a HTML element")
			return
		}

		siteGraph.DOM = *element

		node, err := models.Serialize(siteGraph)

		if err != nil {
			error.Set(Fmt("cannot create site: %v", err))
			return
		}

		if err := node.SaveTree(dbf()); err != nil {
			error.Set(Fmt("cannot save tree: %v", err))
			return
		}

		site.HeadID = &node.ID

		if err := site.Save(); err != nil {
			error.Set("cannot save site")
			return
		}

		fmt.Println(router.CurrentPath())

		router.RedirectTo(router.CurrentPath())

	}

	form.OnSubmit(onSubmit)

	return form.Form(
		If(error.Get() != "", P(error.Get())),
		Textarea(
			Attrib("rows")("20"),
			Styles(Width(Px(600))),
			Value(source),
		),
		Button(
			Type("submit"),
			"Update",
		),
	)
}
