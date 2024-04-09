package ui

import (
	. "github.com/gospel-sh/gospel"
)

func NotFound(c Context) Element {

	c.SetStatusCode(404)

	return Div(
		Class("kip-with-app-selector"),
		A(
			Class("kip-with-app-selector-link"),
			Href("/#"),
			Div(
				Class("kip-logo-wrapper"),
				Img(
					Class("kip-logo", Alt("projects")),
					Src("/static/images/kodexlogo-white.png"),
				),
			),
		),
		Section(
			Class("kip-centered-card", "kip-is-info", "kip-is-fullheight"),
			Div(
				Class("kip-card", "kip-is-centered", "kip-account"),
				Div(
					Class("kip-card-header"),
					Div(
						Class("kip-card-title"),
						H2("404 - Page Not Found"),
					),
				),
				Div(
					Class("kip-card-content", "kip-card-centered"),
					Div(
						"Sorry, there's nothing here...",
					),
				),
			),
		),
	)
}
