package auth

import (
	"fmt"
	"net/http"
)

type SimpleUserProfileProvider struct {
	settings *SimpleSettings
}

type SimpleSettings struct {
	Users []*BasicUserProfile `json:"users"`
}

func (s *SimpleUserProfileProvider) GetWithToken(token []byte) (UserProfile, error) {

	for _, user := range s.settings.Users {
		if string(user.AccessToken().Token()) == string(token) {
			return user, nil
		}
	}
	return nil, fmt.Errorf("invalid access token")
}

func (s *SimpleUserProfileProvider) Get(r *http.Request) (UserProfile, error) {

	token, _, err := GetTokenValue(r)

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, fmt.Errorf("token missing")
	}

	return s.GetWithToken(token)
}

func (s *SimpleUserProfileProvider) Start() {

}

func (s *SimpleUserProfileProvider) Stop() {

}

func (s *SimpleUserProfileProvider) GetWithPassword(email string, password string) (UserProfile, error) {

	for _, user := range s.settings.Users {
		if user.EMail() == email {
			if user.Password() == password {
				return user, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid user or password")
}

func MakeSimpleUserProfileProvider(settings *SimpleSettings) (UserProfileProvider, error) {
	return &SimpleUserProfileProvider{settings: settings}, nil
}
