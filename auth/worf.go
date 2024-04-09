package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/getworf/worf-go"
	"net/http"
	"sync"
	"time"
)

type WorfUserProfileProvider struct {
	worfURL      string
	mutex        *sync.Mutex
	expiresAfter int64
	profiles     map[string]*ProfileWithRetrievalTime
}

type WorfUserProfile struct {
	email       string
	superUser   bool
	displayName string
	sourceID    []byte
	accessToken *WorfAccessToken
	limits      map[string]interface{}
	roles       []OrganizationRoles
}

type WorfOrganizationRoles struct {
	roles        []string
	organization *WorfOrganization
}

type WorfAccessToken struct {
	scopes []string
	token  []byte
}

func MakeWorfUserProfile(profile *worf.UserProfile) *WorfUserProfile {
	orgRoles := make([]OrganizationRoles, 0)

	tokenValue, _ := hex.DecodeString(profile.AccessToken.Token)

	accessToken := &WorfAccessToken{
		scopes: profile.AccessToken.Scopes,
		token:  tokenValue,
	}

	if len(profile.Organizations) > 0 {
		roles := make([]string, 0)
		for _, organization := range profile.Organizations {
			for _, role := range organization.Roles {
				if !role.Confirmed {
					continue
				}
				roles = append(roles, role.Role)
			}

			organization := &WorfOrganization{
				source:      "worf",
				id:          organization.BinaryID(),
				description: organization.Description,
				name:        organization.Name,
			}

			or := &WorfOrganizationRoles{
				roles:        roles,
				organization: organization,
			}

			orgRoles = append(orgRoles, or)

		}

	}

	userOrgName := profile.User.DisplayName

	if userOrgName == "" {
		userOrgName = profile.User.EMail
	}

	userOrganization := &WorfOrganization{
		source:   "worf_user",
		id:       profile.User.BinaryID(),
		name:     userOrgName,
		default_: true,
	}

	sorg := &WorfOrganizationRoles{
		roles:        []string{"superuser", "admin"},
		organization: userOrganization,
	}

	orgRoles = append(orgRoles, sorg)

	userProfile := &WorfUserProfile{
		sourceID:    profile.User.BinaryID(),
		email:       profile.User.EMail,
		superUser:   profile.User.SuperUser,
		displayName: profile.User.DisplayName,
		accessToken: accessToken,
		limits:      profile.Limits,
		roles:       orgRoles,
	}

	return userProfile
}

func (w *WorfUserProfile) Source() string {
	return "worf"
}

func (w *WorfUserProfile) SourceID() []byte {
	return w.sourceID
}

func (w *WorfUserProfile) EMail() string {
	return w.email
}

func (w *WorfUserProfile) SuperUser() bool {
	return w.superUser
}

func (w *WorfUserProfile) Limits() map[string]interface{} {
	return w.limits
}

func (w *WorfUserProfile) DisplayName() string {
	return w.displayName
}

func (w *WorfUserProfile) AccessToken() AccessToken {
	return w.accessToken
}

func (w *WorfUserProfile) Roles() []OrganizationRoles {
	return w.roles
}

func (w *WorfAccessToken) Scopes() []string {
	return w.scopes
}

func (w *WorfAccessToken) Token() []byte {
	return w.token
}

func (w *WorfOrganizationRoles) Roles() []string {
	return w.roles
}

func (w *WorfOrganizationRoles) Organization() UserOrganization {
	return w.organization
}

type WorfOrganization struct {
	name        string
	source      string
	description string
	default_    bool
	id          []byte
}

func (w *WorfOrganization) Default() bool {
	return w.default_
}

func (w *WorfOrganization) Name() string {
	return w.name
}

func (w *WorfOrganization) Source() string {
	return w.source
}

func (w *WorfOrganization) Description() string {
	return w.description
}

func (w *WorfOrganization) ID() []byte {
	return w.id
}

func MakeWorfUserProfileProvider(settings map[string]any) (UserProfileProvider, error) {
	worfApiURL, ok := settings["url"].(string)

	if !ok {
		return nil, fmt.Errorf("Worf URL missing (worf.url)")
	}

	expiresAfter, ok := settings["expires-after"].(int64)

	if !ok {
		expiresAfter = 60
	}

	return &WorfUserProfileProvider{
		profiles:     make(map[string]*ProfileWithRetrievalTime),
		worfURL:      worfApiURL,
		expiresAfter: expiresAfter,
		mutex:        &sync.Mutex{},
	}, nil
}

func getProfileFromAPI(apiURL string, accessToken []byte) (UserProfile, error) {
	client := worf.MakeClient(apiURL, hex.EncodeToString(accessToken))
	worfProfile, err := client.UserProfile()
	if err != nil {
		return nil, err
	}
	return MakeWorfUserProfile(worfProfile), nil
}

func hash(accessToken []byte) string {
	h := sha256.Sum256(accessToken)
	hb := h[:]
	return string(hb)
}

func (a *WorfUserProfileProvider) refreshProfile(token []byte) {

	userProfile, err := getProfileFromAPI(a.worfURL, token)

	h := hash(token)
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err != nil {
		delete(a.profiles, h)
		return
	}

	a.profiles[h] = &ProfileWithRetrievalTime{
		Profile:     userProfile,
		RetrievedAt: time.Now().UTC(),
		Refreshing:  false,
		mutex:       &sync.Mutex{},
	}

}

func (a *WorfUserProfileProvider) Start() {
	go a.cleanProfiles()
}

func (a *WorfUserProfileProvider) Stop() {

}

func (a *WorfUserProfileProvider) cleanProfiles() {
	for {
		time.Sleep(time.Second * 1)
		a.mutex.Lock()
		for key, profile := range a.profiles {
			if profile.RetrievedAt.Add(time.Second * time.Duration(a.expiresAfter)).Before(time.Now().UTC()) {
				delete(a.profiles, key)
			}
		}
		a.mutex.Unlock()
	}
}

func GetTokenValue(r *http.Request) ([]byte, bool, error) {

	if cookie, err := r.Cookie("auth"); err == nil {
		if token, err := hex.DecodeString(cookie.Value); err == nil && len(token) >= 16 {
			return token, true, nil
		}
	}

	token := r.Header.Get("Authorization")

	if token == "" || len(token) < 8 {
		return nil, false, nil
	}

	// we remove the 'Bearer ' part...
	token = token[7:]

	if tokenBytes, err := hex.DecodeString(token); err != nil {
		return nil, false, err
	} else {
		return tokenBytes, false, nil
	}
}

func (a *WorfUserProfileProvider) GetWithPassword(email, password string) (UserProfile, error) {

	client := worf.MakeClient(a.worfURL, "")
	response, err := client.PasswordLogin(email, password)

	if err != nil {
		return nil, err
	}

	return MakeWorfUserProfile(response), nil
}

func (a *WorfUserProfileProvider) GetWithToken(token []byte) (UserProfile, error) {
	a.mutex.Lock()
	h := hash(token)
	profile, ok := a.profiles[h]
	a.mutex.Unlock()

	if ok {
		profile.mutex.Lock()
		defer profile.mutex.Unlock()
		if (!profile.Refreshing) && profile.RetrievedAt.Add(time.Second*time.Duration(a.expiresAfter)/2).Before(time.Now().UTC()) {
			profile.Refreshing = true
			go a.refreshProfile(token)
		}
		return profile.Profile, profile.Error
	}

	userProfile, err := getProfileFromAPI(a.worfURL, token)

	if err != nil {
		return nil, err
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.profiles[h] = &ProfileWithRetrievalTime{
		Profile:     userProfile,
		Error:       err,
		RetrievedAt: time.Now().UTC(),
		mutex:       &sync.Mutex{},
	}

	return userProfile, err
}

func (a *WorfUserProfileProvider) Get(r *http.Request) (UserProfile, error) {

	token, _, err := GetTokenValue(r)

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, fmt.Errorf("token missing")
	}

	return a.GetWithToken(token)

}
