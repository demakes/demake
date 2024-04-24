package ui

import (
	"encoding/hex"
	"fmt"
	. "github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites/auth"
	"github.com/klaro-org/sites/models"
)

func SetDB(c Context, db orm.DB) {
	GlobalVar(c, "db", db)
}

func UseDB(c Context) orm.DB {
	return UseGlobal[orm.DB](c, "db")
}

func SetUser(c Context, user auth.UserProfile) {
	GlobalVar(c, "user", user)
}

func UseUser(c Context) auth.UserProfile {
	return UseGlobal[auth.UserProfile](c, "user")
}

func SetProfileProvider(c Context, provider auth.UserProfileProvider) {
	GlobalVar(c, "profileProvider", provider)
}

func UseProfileProvider(c Context) auth.UserProfileProvider {
	return UseGlobal[auth.UserProfileProvider](c, "profileProvider")
}

var textFont = FontFamily("'Poppins', sans-serif")
var titleFont = FontFamily("'Bricolage Grotesque', sans-serif")

var CSS = MakeStylesheet(
	"root",
	Html(
		FontSize(Px(18)),
		Height("100%"),
		Margin(0),
		Padding(0),
	),
	Body(
		Height("100%"),
		Margin(0),
		Padding(0),
		H1(
			titleFont,
			FontSize(Rem(2.4)),
			FontWeight(300),
		),
		H2(
			titleFont,
			FontSize(Rem(2.0)),
			FontWeight(300),
		),
		H3(
			titleFont,
			FontSize(Rem(1.6)),
			FontWeight(300),
		),
		Any(
			textFont,
			FontSize(Rem(1.0)),
			FontWeight(400),
		),
	),
)

func MainContent(c Context) Element {

	router := UseRouter(c)

	// if the user isn't logged in, we redirect to the login screen
	if UseUser(c) == nil {
		router.RedirectTo("/login")
		return nil
	}

	return AuthorizedContent(c)

}

func AuthorizedContent(c Context) Element {

	router := UseRouter(c)

	return Div(
		Styles(
			Display("flex"),
		),
		MainHeader(c),
		Div(
			Styles(
				MarginTop(Px(80)),
			),
			router.Match(
				c,
				Route(
					"/sites",
					Sites,
				),
				Route(
					"",
					NotFound,
				),
			),
		),
	)
}

func GetGraph(site *models.Site, dbf func() orm.DB) (*models.SiteGraph, error) {
	if site.HeadID == nil {
		return nil, fmt.Errorf("site doesn't have a head")
	}

	graph, err := models.GetGraphByID(dbf, *site.HeadID)

	if err != nil {
		return nil, fmt.Errorf("cannot get graph: %v", err)
	}

	siteGraph, err := models.DeserializeType[models.SiteGraph](graph)

	if err != nil {
		return nil, fmt.Errorf("cannot deserialize graph: %v", err)
	}

	return siteGraph, nil

}

func ServeSite(db orm.DB, site *models.Site) func(c Context) Element {

	dbf := func() orm.DB { return db }

	return func(c Context) Element {

		siteGraph, err := GetGraph(site, dbf)

		if err != nil {
			return Div(Fmt("Cannot load site: %v", err))
		}

		element, err := siteGraph.DOM.Generate(c)

		if err != nil {
			return Div("cannot generate")
		}

		return element.(Element)
	}

}

func Root(db orm.DB, profileProvider auth.UserProfileProvider) func(c Context) Element {

	dbf := func() orm.DB { return db }

	return func(c Context) Element {

		SetDB(c, db)
		SetProfileProvider(c, profileProvider)

		// if the user isn't logged in, we redirect to the login screen
		if user, err := profileProvider.Get(c.Request()); err == nil {
			SetUser(c, user)
		}

		router := UseRouter(c)
		router.SetPrefix("/demake")

		site := router.Match(
			c,
			Route(`/sites/view/([a-f0-9\-]+)`, func(c Context, siteID string) Element {

				id, err := hex.DecodeString(siteID)

				if err != nil {
					return Div("invalid ID")
				}

				if UseUser(c) == nil {
					router.RedirectTo("/login")
					return nil
				}

				site := orm.Init(&models.Site{}, dbf)

				if err := site.ByExtID(id); err != nil {
					return Div("cannot find site")
				}

				return ServeSite(db, site)(c)

			}),
		)

		if site != nil {
			return site
		}

		return F(
			Doctype("html"),
			Html(
				Lang("en"),
				Head(
					Meta(Charset("utf-8")),
					Title("Demake"),
					CSS.Styles(),
					// to do: remove this
					L(`<link rel="preconnect" href="https://fonts.googleapis.com"><link rel="preconnect" href="https://fonts.gstatic.com" crossorigin><link href="https://fonts.googleapis.com/css2?family=Bricolage+Grotesque:opsz,wght@12..96,300;12..96,600;12..96,800&family=Darker+Grotesque:wght@300;700;900&family=Poppins:wght@100;400;700&family=Roboto:wght@100;300;700&display=swap" rel="stylesheet">`),
				),
				Body(
					Styled(
						"main",
						router.Match(
							c,
							Route("/login", Login),
							Route("/logout", Logout),
							Route("/404", NotFound),
							Route("", MainContent),
						),
					),
				),
			),
		)
	}

}
