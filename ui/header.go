package ui

import (
	. "github.com/gospel-sh/gospel"
)

type Breadcrumb struct {
	Title string
	Path  string
}

func Breadcrumbs(c Context) Element {

	crumbs := []Element{}

	breadcrumbs := UseGlobal[[]Breadcrumb](c, "breadcrumbs")

	path := ""
	title := ""

	for _, breadcrumb := range breadcrumbs {

		path += breadcrumb.Path

		if title != "" {
			title += " :: "
		}

		title += breadcrumb.Title

		crumbs = append(crumbs, Li(
			A(Href(path), breadcrumb.Title),
		))
	}

	return Nav(
		Class("bulma-breadcrumb bulma-has-bullet-separator"),
		Ul(
			crumbs,
		),
	)

}

func AddBreadcrumb(c Context, title string, path string) {

	breadcrumbs := GlobalVar(c, "breadcrumbs", []Breadcrumb{})

	bcs := breadcrumbs.Get()

	bcs = append(bcs, Breadcrumb{
		Title: title,
		Path:  path,
	})

	breadcrumbs.Set(bcs)
}

func MainTitle(c Context) string {

	breadcrumbs := GlobalVar(c, "breadcrumbs", []Breadcrumb{})

	title := ""

	for _, breadcrumb := range breadcrumbs.Get() {

		if title != "" {
			title += " :: "
		}

		title += breadcrumb.Title
	}

	// we reset the breadcrumbs
	breadcrumbs.Set([]Breadcrumb{})

	return title

}

func MainHeader(c Context) Element {

	user := UseUser(c)
	breadcrumbs := UseGlobal[[]Breadcrumb](c, "breadcrumbs")
	tab := ""

	if len(breadcrumbs) > 0 {
		tab = breadcrumbs[0].Path
	}

	return Header(
		Nav(
			Styles(
				Position("absolute"),
				Display("flex"),
				AlignItems("center"),
				Top(0),
				Left(0),
				Background("#eee"),
				PaddingLeft(Px(20)),
				Width(Calc(Sub(Percent(100), Px(20)))),
				Display("flex"),
				Margin(0),
				BorderBottom("2px solid #333"),
			),
			Ul(
				Styles(
					Margin(0),
					FlexGrow(1),
					Color("#333"),
					ListStyle("none"),
					Display("flex"),
					PaddingLeft(Px(20)),
					FontWeight(600),
					Direct(Li)(
						Display("inline-block"),
						MarginRight(Px(10)),
						Display("flex"),
						JustifyContent("end"),
						LastChild(
							MarginRight(0),
						),
						ApplyAny(A, Span)(
							Class("is-active",
								Background("#fafafa"),
							),
							Hover(
								Background("#ddd"),
							),
							Color("#333"),
							Padding(Px(6)),
							Margin("auto 0"),
							TextDecoration("none"),
						),
					),
				),
				Li(A(Href("/sites"), "Sites", If(tab == "sites", Class("is-active")))),
				Li(
					Styles(
						FlexGrow(1),
						Display("flex"),
					),
					Span(
						Styles(
							Position("relative"),
							Ul(
								Display("none"),
							),
							Hover(
								Background("#eee"),
								Ul(
									Margin(0),
									PaddingLeft(0),
									boxShadow(6),
									Border("2px solid #000"),
									Display("flex"),
									FlexDirection("column"),
									ListStyle("none"),
									Width(Percent(90)),
									Background("#eee"),
									Position("absolute"),
									Left(0),
									Li(
										Padding(0),
										Display("flex"),
										Margin(0),
										A(
											FlexGrow(1),
										),
									),
								),
							),
						),
						Span(user.EMail()),
						Ul(
							Li(
								A(Href("/logout"), "Logout"),
							),
						),
					),
				),
			),
		),
	)
}
