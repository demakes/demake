package sites

import (
	"fmt"
)

// NotFound is returned when a given object can't be found
var NotFound = fmt.Errorf("not found")

// A Backend can retrieve objects
type Backend interface {
	// Sites

	// Retrieve all sites matching the given filters
	Sites(filters map[string]any) ([]*Site, error)
}

// A WritableBackend can write and delete objects
type WritableBackend interface {
	// Sites
	SaveSite(site *Site) error
	DeleteSite(site *Site) error
}
