package ui

import (
	. "github.com/gospel-sh/gospel"
	"github.com/klaro-org/sites/auth"
	"net/http"
	"time"
)

func boxShadow(size int) *Declaration {
	return BoxShadow(Fmt("rgba(0, 0, 0, 0.5) 0px 0px 0px 0px, rgba(0, 0, 0, 0.5) 0px 0px 0px 0px, rgba(0, 0, 0, 0.5) %[2]dpx %[1]dpx 0px 0px", size, size*2/3))
}

func shadowed(radius float64, size int) []any {
	return []any{
		boxShadow(size),
		BorderRadius(Px(radius)),
		Border("2px solid #000"),
		Padding(Px(20)),
	}
}

func shadowedButton(radius float64, size int) []any {
	return append(shadowed(radius, size), []any{
		Active(
			BoxShadow("none"),
			Transition("0s"),
			Transform(Fmt("translate(%[1]dpx, %[1]dpx)", size)),
		),
		Transform("scale(100%)"),
		Hover(
			Filter("brightness(110%)"),
		),
	}...)
}

func Login(c Context) Element {

	form := MakeFormData(c, "login", POST)
	email := form.Var("email", "")
	password := form.Var("password", "")
	error := Var(c, "")
	router := UseRouter(c)
	profileProvider := UseProfileProvider(c)
	passwordProvider, ok := profileProvider.(auth.PasswordProvider)

	if !ok {
		return Div("cannot log in via password")
	}

	onSubmit := Func[any](c, func() {

		if email.Get() == "" {
			error.Set("Please enter an e-mail")
			return
		}

		if password.Get() == "" {
			error.Set("Please enter a password")
			return
		}

		// we check the token
		if profile, err := passwordProvider.GetWithPassword(email.Get(), password.Get()); err != nil {
			error.Set("invalid password or username")
			return
		} else {
			w := c.ResponseWriter()
			http.SetCookie(w, &http.Cookie{Path: "/", Name: "auth", Value: Hex(profile.AccessToken().Token()), Secure: false, HttpOnly: true, Expires: time.Now().Add(365 * 24 * 7 * time.Hour)})
			router.RedirectTo("/")
		}
	})

	return Section(
		// Background
		Styles(
			Width(Percent(100)),
			Height(Percent(100)),
			Display("flex"),
			FlexDirection("row"),
			Background("#aaa"),
			AlignItems("center"),
			JustifyContent("center"),
		),
		Div(
			// Card
			Styles(
				Width(Px(500)),
				Height("auto"),
				Background("#C6DDE2"),
				shadowed(20, 20),
			),
			// Header
			Div(
				Styles(
					Padding(Px(20)),
				),
				Div(
					Styles(
						H2(
							TextAlign("center"),
							MarginTop(Px(6)),
							FontWeight(400),
							TextDecoration("underline"),
							MarginBottom(Px(6)),
						),
					),
					H2("Login"),
				),
			),
			// Content
			Div(
				Styles(
					Padding(Px(40)),
					PaddingTop(0),
				),
				Form(
					Styles(
						Display("flex"),
						FlexDirection("column"),
						Label(
							Display("block"),
							MarginTop(Rem(1.0)),
						),
						Input(
							shadowed(6, 6),
							Background("#eee"),
							Width(Percent(100)),
							BoxSizing("border-box"),
							If(
								error.Get() != "",
								[]any{
									BorderColor("#a66"),
									Color("#a66"),
								},
							),
						),
					),
					Method("POST"),
					OnSubmit(onSubmit),
					Div(
						Styles(
							FlexGrow(1),
						),
						P(
							Styles(
								Color("#a66"),
							),
							IfElse(error.Get() != "", Span(error.Get()), Nbsp),
						),
						Label(
							"E-Mail",
							Input(
								Value(email),
								Placeholder("e-mail"),
								Type("email"),
							),
						),
						Label(
							"Password",
							Input(
								Value(password),
								Type("password"),
								Placeholder("password"),
							),
						),
					),
					Div(
						P(
							Button(
								Styles(
									shadowedButton(6, 6),
									Background("#efa"),
								),
								Type("submit"),
								"Log in",
							),
						),
					),
				),
			),
		),
	)
}
