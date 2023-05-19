package klaro

// A Route links to a page
type Route struct {
	Base
	Path   string `json:"path"`
	PageID uint64 `json:"pageId"`
	SiteID uint64 `json:"siteId"`
	// A cached version of the associated page
	page *Page `json:"-"`
}

func (r *Route) Page() (*Page, error) {
	if r.page == nil {

		var err error

		if r.page, err = r.Backend.Page(r.PageID); err != nil {
			return nil, err
		}
	}

	return r.page, nil
}
