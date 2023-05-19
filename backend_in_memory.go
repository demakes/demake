package klaro

type InMemoryBackend struct {
	sites  []*Site
	routes []*Route
	pages  []*Page
}

func MakeInMemoryBackend() *InMemoryBackend {
	return &InMemoryBackend{
		sites:  make([]*Site, 0),
		routes: make([]*Route, 0),
		pages:  make([]*Page, 0),
	}
}

func (b *InMemoryBackend) Sites() ([]*Site, error) {
	return b.sites, nil
}

func (b *InMemoryBackend) Routes(siteId uint64) ([]*Route, error) {
	routes := make([]*Route, 0)

	for _, route := range b.routes {
		if route.SiteID == siteId {
			routes = append(routes, route)
		}
	}
	return routes, nil
}

func (b *InMemoryBackend) Pages(siteId uint64, filters map[string]any) ([]*Page, error) {
	pages := make([]*Page, 0)

	for _, page := range b.pages {
		if page.SiteID == siteId {
			pages = append(pages, page)
		}
	}

	return pages, nil
}

func (b *InMemoryBackend) Page(pageId uint64) (*Page, error) {

	for _, page := range b.pages {
		if page.ID == pageId {
			return page, nil
		}
	}
	return nil, NotFound
}
