package klaro

// A Page produces a webpage
type Page struct {
	Base
	SiteID uint64 `json:"siteId"`
}
