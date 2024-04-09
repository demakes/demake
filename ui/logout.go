package ui

import (
	. "github.com/gospel-sh/gospel"
	"net/http"
	"time"
)

func Logout(c Context) Element {

	w := c.ResponseWriter()

	http.SetCookie(w, &http.Cookie{Path: "/", Name: "auth", Value: "", Secure: false, HttpOnly: true, Expires: time.Unix(0, 0)})

	// we clear the context
	c.Clear()

	return Section(
		Class("kip-centered-card", "kip-is-info", "kip-is-fullheight"),
		Div(
			Class("kip-card", "kip-is-centered", "kip-account"),
			Div(
				Class("kip-card-header"),
				Div(
					Class("kip-card-title"),
					H2("Logout"),
				),
			),
			Div(
				Class("kip-card-content", "kip-card-centered"),
				P(
					"You have been logged out. ",
					A(Href("/login"), "Log back in."),
				),
			),
		),
	)
}
