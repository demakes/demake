package ui

import (
	. "github.com/gospel-sh/gospel"
)

func MainMenu(c Context) Element {
	return Aside(
		// Menu bar
		Styles(
			Width(Px(300)),
			Position("fixed"),
			Top(Px(0)),
			Left(0),
			Height(Percent(100)),
			Background("#D6CEFD"),
			BorderRight("2px solid #333"),
			Mobile(
				Display("none"),
			),
		),
		Div(
			Styles(
				Padding(Px(20)),
			),
			Div(
				Styles(
					Color("#000"),
					Display("flex"),
					MarginBottom(Rem(2)),
					AlignItems("center"),
					Svg(
						Height(Rem(2)),
						MarginRight(Rem(1)),
					),
				),
				Logo(),
				Span(
					Styles(),
					"TIMELIKE",
				),
			),
			Ul(
				Styles(
					Padding(0),
					ListStyle("none"),
					Margin(0),
					MarginTop(Px(10)),
					MarginLeft(Px(-6)),
					Li(
						Margin(Px(8)),
						Padding(Px(6)),
						MarginLeft(0),
					),
					A(
						Color("#000"),
						Background("#aae"),
						Display("flex"),
						AlignItems("center"),
						shadowedButton(0, 6),
						TextDecoration("none"),
						Svg(
							Stroke("#000"),
							StrokeWidth(1.5),
							Fill("none"),
							FlexShrink(0),
							Width(Rem(1.5)),
							Height(Rem(1.5)),
							MarginRight(Px(10)),
						),
					),
				),
				Li(
					A(
						Href("#"),
						Svg(
							L(`<path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"></path>`),
						),
						"Dashboard",
					),
				),
				Li(
					A(
						Href("#"),
						Svg(
							L(`<path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"></path>`),
						),
						"Users",
					),
				),
				Li(
					A(
						Href("#"),
						Svg(
							L(`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6"><path strokeLinecap="round" strokeLinejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5m-9-6h.008v.008H12v-.008zM12 15h.008v.008H12V15zm0 2.25h.008v.008H12v-.008zM9.75 15h.008v.008H9.75V15zm0 2.25h.008v.008H9.75v-.008zM7.5 15h.008v.008H7.5V15zm0 2.25h.008v.008H7.5v-.008zm6.75-4.5h.008v.008h-.008v-.008zm0 2.25h.008v.008h-.008V15zm0 2.25h.008v.008h-.008v-.008zm2.25-4.5h.008v.008H16.5v-.008zm0 2.25h.008v.008H16.5V15z" /></svg>`),
						),
						"Schedules",
					),
				),
			),
		),
	)
}
