package ui

import (
	. "github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites/models"
)

func getSites(c Context) ([]*models.Site, error) {
	db := func() orm.DB { return UseDB(c) }
	return orm.Objects[models.Site](db, map[string]any{})
}

func NewSite(c Context) Element {

	db := func() orm.DB { return UseDB(c) }

	error := Var[string](c, "")
	name := Var[string](c, "")
	onSubmit := Func[any](c, func() {

		if len(name.Get()) == 0 {
			error.Set("please enter a name")
			return
		}

		sites, err := getSites(c)

		if err != nil {
			error.Set("cannot load sites")
			return
		}

		for _, site := range sites {
			if site.Name == name.Get() {
				error.Set("a site with this name already exists")
				return
			}
		}

		newSite := &models.Site{
			Name: name.Get(),
		}

		orm.Init(newSite, db)

		if err := newSite.Save(); err != nil {
			error.Set("cannot save site")
			return
		}

		UseRouter(c).RedirectTo("/sites")

	})

	return Div(
		Form(
			Method("POST"),
			OnSubmit(onSubmit),
			If(error.Get() != "", error.Get()),
			Input(Value(name)),
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
			site.Name,
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
		A(Href("/sites/new"), "new site"),
	)
}

func Sites(c Context) Element {

	AddBreadcrumb(c, "Sites", "sites")

	return Div(
		UseRouter(c).Match(
			c,
			Route("/new$", NewSite),
			Route("$", SiteList),
		),
	)

}
