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

// Implementation

type BasicUserProfileFields struct {
	EMail       string                 `json:"email"`
	SuperUser   bool                   `json:"superuser"`
	DisplayName string                 `json:"displayname"`
	Source      string                 `json:"source"`
	SourceID    []byte                 `json:"sourceID"`
	AccessToken *BasicAccessToken      `json:"accessToken"`
	Limits      map[string]interface{} `json:"limits"`
	Roles       []OrganizationRoles    `json:"roles"`
	// only used for the simple provider
	Password string `json:"password"`
}

type BasicUserProfile struct {
	BasicUserProfileFields
}

type BasicOrganizationRolesFields struct {
	Roles        []string           `json:"roles"`
	Organization *BasicOrganization `json:"organization"`
}

type BasicOrganizationRoles struct {
	BasicOrganizationRolesFields
}

type BasicAccessTokenFields struct {
	Scopes []string `json:"scopes"`
	Token  []byte   `json:"token"`
}

type BasicAccessToken struct {
	BasicAccessTokenFields
}

type BasicOrganizationFields struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
	ID          []byte `json:"id"`
}

type BasicOrganization struct {
	BasicOrganizationFields
}

func (w *BasicUserProfile) Password() string {
	return w.BasicUserProfileFields.Password
}

func (w *BasicUserProfile) Source() string {
	return w.BasicUserProfileFields.Source
}

func (w *BasicUserProfile) SourceID() []byte {
	return w.BasicUserProfileFields.SourceID
}

func (w *BasicUserProfile) EMail() string {
	return w.BasicUserProfileFields.EMail
}

func (w *BasicUserProfile) SuperUser() bool {
	return w.BasicUserProfileFields.SuperUser
}

func (w *BasicUserProfile) Limits() map[string]interface{} {
	return w.BasicUserProfileFields.Limits
}

func (w *BasicUserProfile) DisplayName() string {
	return w.BasicUserProfileFields.DisplayName
}

func (w *BasicUserProfile) AccessToken() AccessToken {
	return w.BasicUserProfileFields.AccessToken
}

func (w *BasicUserProfile) Roles() []OrganizationRoles {
	return w.BasicUserProfileFields.Roles
}

func (w *BasicAccessToken) Scopes() []string {
	return w.BasicAccessTokenFields.Scopes
}

func (w *BasicAccessToken) Token() []byte {
	return w.BasicAccessTokenFields.Token
}

func (w *BasicOrganizationRoles) Roles() []string {
	return w.BasicOrganizationRolesFields.Roles
}

func (w *BasicOrganizationRoles) Organization() UserOrganization {
	return w.BasicOrganizationRolesFields.Organization
}

func (w *BasicOrganization) Default() bool {
	return w.BasicOrganizationFields.Default
}

func (w *BasicOrganization) Name() string {
	return w.BasicOrganizationFields.Name
}

func (w *BasicOrganization) Source() string {
	return w.BasicOrganizationFields.Source
}

func (w *BasicOrganization) Description() string {
	return w.BasicOrganizationFields.Description
}

func (w *BasicOrganization) ID() []byte {
	return w.BasicOrganizationFields.ID
}
