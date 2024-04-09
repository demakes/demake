package ui

import (
	. "github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites/auth"
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

	profileProvider := UseProfileProvider(c)
	router := UseRouter(c)

	// if the user isn't logged in, we redirect to the login screen
	if user, err := profileProvider.Get(c.Request()); err != nil {
		router.RedirectTo("/login")
		return nil
	} else {
		SetUser(c, user)
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

func Root(db orm.DB, profileProvider auth.UserProfileProvider) func(c Context) Element {

	return func(c Context) Element {

		SetDB(c, db)
		SetProfileProvider(c, profileProvider)

		router := UseRouter(c)

		return F(
			Doctype("html"),
			Html(
				Lang("en"),
				Head(
					Meta(Charset("utf-8")),
					Title("Linearize"),
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
