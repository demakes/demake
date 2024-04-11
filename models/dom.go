package models

import (
	"github.com/gospel-sh/gospel"
)

func init() {
	MustRegister[gospel.HTMLElement]("element")
	MustRegister[gospel.HTMLAttribute]("attribute")
	MustRegister[gospel.RouteConfig]("route")
}
