package klaro

import (
	"fmt"
)

var NotFound = fmt.Errorf("not found")

type Backend interface {
	Sites() ([]*Site, error)
	Routes(siteId uint64) ([]*Route, error)
	Pages(siteId uint64, filters map[string]any) ([]*Page, error)
	Page(pageId uint64) (*Page, error)
}
