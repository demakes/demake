package ui

import (
	"github.com/demakes/demake/models"
	. "github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
)

func getSites(c Context) ([]*models.Site, error) {
	db := func() orm.DB { return UseDB(c) }
	return orm.Objects[models.Site](db, map[string]any{})
}

func NewSite(c Context) Element {

	db := func() orm.DB { return UseDB(c) }

	formData := MakeFormData(c, "newSite", POST)

	error := Var[string](c, "")
	name := formData.Var("name", "")
	hostname := formData.Var("hostname", "")
	onSubmit := func() {

		if len(name.Get()) == 0 {
			error.Set("please enter a name")
			return
		}

		if len(name.Get()) == 0 {
			error.Set("please enter a hostname")
			return
		}

		sites, err := getSites(c)

		if err != nil {
			error.Set(Fmt("cannot load sites: %v", err))
			return
		}

		for _, site := range sites {
			if site.Hostname == hostname.Get() || site.Name == name.Get() {
				error.Set("a site with this name or hostname already exists")
				return
			}
		}

		newSite := &models.Site{
			Name:     name.Get(),
			Hostname: hostname.Get(),
		}

		orm.Init(newSite, db)

		dom := Div(
			P("This is a test"),
			Strong("strongs"),
			Route("/test(/[a-z]+)?", Strong("another test")),
			Route("/blub(/[a-z]+)?", Strong("blub")),
			Route("/foo(/[a-z]+)?", Strong("foo")),
		)

		site := &models.SiteGraph{
			Meta: models.SiteMeta{
				Title: models.TranslatedString{
					Translations: map[string]string{"de": "Meine Webseite"},
				},
				Domain: "japh.de",
			},
			DOM:     *dom,
			Plugins: []models.SitePlugin{&models.BlogPlugin{ArticlesPerPage: 10}},
		}

		node, err := models.Serialize(site)

		if err != nil {
			error.Set(Fmt("cannot create site: %v", err))
			return
		}

		if err := node.SaveTree(db()); err != nil {
			error.Set(Fmt("cannot save tree: %v", err))
			return
		}

		newSite.HeadID = &node.ID

		if err := newSite.Save(); err != nil {
			error.Set("cannot save site")
			return
		}

		UseRouter(c).RedirectTo("/sites")
	}

	formData.OnSubmit(onSubmit)

	return Div(
		formData.Form(
			If(error.Get() != "", error.Get()),
			Input(Placeholder("name"), Value(name)),
			Input(Placeholder("hostname"), Value(hostname)),
			Button(
				Type("submit"),
				"create site",
			),
		),
	)

}

func SiteList(c Context) Element {

	sites, err := getSites(c)

	if err != nil {
		return Div(Fmt("error: %v", err))
	}

	siteItems := make([]Element, len(sites))

	for i, site := range sites {
		siteItems[i] = Li(
			A(
				Href(UseRouter(c).URL(Fmt("/sites/edit/%s", site.ExtID.Hex()))),
				site.Name,
			),
			" // ",
			site.CreatedAt.String(),
			" // ",
			site.UpdatedAt.String(),
		)
	}

	return Div(
		Ul(
			siteItems,
		),
		A(Href(UseRouter(c).URL("/sites/new")), "new site"),
	)
}

func Sites(c Context) Element {

	AddBreadcrumb(c, "Sites", "sites")

	return Div(
		UseRouter(c).Match(
			c,
			Route("/new$", NewSite),
			Route(`/edit/([a-f0-9\-]+)`, EditSite),
			Route("$", SiteList),
		),
	)

}
