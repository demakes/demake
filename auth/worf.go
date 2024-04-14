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

type WorfSettings struct {
	URL          string `json:"url"`
	ExpiresAfter int64  `json:"expiresAfter"`
}

func MakeWorfUserProfile(profile *worf.UserProfile) *BasicUserProfile {
	orgRoles := make([]OrganizationRoles, 0)

	tokenValue, _ := hex.DecodeString(profile.AccessToken.Token)

	accessToken := &BasicAccessToken{
		BasicAccessTokenFields{
			Scopes: profile.AccessToken.Scopes,
			Token:  tokenValue,
		},
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

			organization := &BasicOrganization{
				BasicOrganizationFields{
					Source:      "worf",
					ID:          organization.BinaryID(),
					Description: organization.Description,
					Name:        organization.Name,
				},
			}

			or := &BasicOrganizationRoles{
				BasicOrganizationRolesFields{
					Roles:        roles,
					Organization: organization,
				},
			}

			orgRoles = append(orgRoles, or)

		}

	}

	userOrgName := profile.User.DisplayName

	if userOrgName == "" {
		userOrgName = profile.User.EMail
	}

	userOrganization := &BasicOrganization{
		BasicOrganizationFields{
			Source:  "worf_user",
			ID:      profile.User.BinaryID(),
			Name:    userOrgName,
			Default: true,
		},
	}

	sorg := &BasicOrganizationRoles{
		BasicOrganizationRolesFields{
			Roles:        []string{"superuser", "admin"},
			Organization: userOrganization,
		},
	}

	orgRoles = append(orgRoles, sorg)

	userProfile := &BasicUserProfile{
		BasicUserProfileFields{
			SourceID:    profile.User.BinaryID(),
			EMail:       profile.User.EMail,
			SuperUser:   profile.User.SuperUser,
			DisplayName: profile.User.DisplayName,
			AccessToken: accessToken,
			Limits:      profile.Limits,
			Roles:       orgRoles,
			Source:      "worf",
		},
	}

	return userProfile
}

func MakeWorfUserProfileProvider(settings *WorfSettings) (UserProfileProvider, error) {

	if settings.ExpiresAfter == 0 {
		settings.ExpiresAfter = 60
	}

	return &WorfUserProfileProvider{
		profiles:     make(map[string]*ProfileWithRetrievalTime),
		worfURL:      settings.URL,
		expiresAfter: settings.ExpiresAfter,
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
		if token, err := hex.DecodeString(cookie.Value); err == nil && len(token) >= 4 {
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
