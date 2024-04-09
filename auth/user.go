package auth

type UserProfile interface {
	Source() string
	SourceID() []byte
	EMail() string
	SuperUser() bool
	DisplayName() string
	AccessToken() AccessToken
	Roles() []OrganizationRoles
	Limits() map[string]interface{}
}

type AccessToken interface {
	Scopes() []string
	Token() []byte
}

type OrganizationRoles interface {
	Roles() []string
	Organization() UserOrganization
}

type UserOrganization interface {
	Name() string
	Source() string
	Default() bool
	Description() string
	ID() []byte
	// ApiOrganization(Controller) (Organization, error)
}
