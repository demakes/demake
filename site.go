package klaro

// A Site contains many routes and pages
type Site struct {
	Base
	routes []*Route `json:"routes"`
}

func (s *Site) Routes() ([]*Route, error) {

	if s.routes == nil {

		var err error

		if s.routes, err = s.Backend.Routes(s.ID); err != nil {
			return nil, err
		}
	}

	return s.routes, nil
}

func (s *Site) Pages(filters map[string]any) ([]*Page, error) {
	return s.Backend.Pages(s.ID, filters)
}
