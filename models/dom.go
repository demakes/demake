package models

import (
	"github.com/gospel-sh/gospel"
)

func init() {
	if err := Register[gospel.HTMLElement]("element"); err != nil {
		panic("cannot register HTMLElement")
	}
	if err := Register[gospel.HTMLAttribute]("attribute"); err != nil {
		panic("cannot register HTMLAttribute")
	}
	if err := Register[gospel.RouteConfig]("route"); err != nil {
		panic("cannot register route config")
	}
}
