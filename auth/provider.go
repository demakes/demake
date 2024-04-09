package auth

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ProfileWithRetrievalTime struct {
	Profile     UserProfile
	Error       error
	RetrievedAt time.Time
	Refreshing  bool
	mutex       *sync.Mutex
}

type UserProfileProviderMaker func(settings map[string]any) (UserProfileProvider, error)

type UserProfileProvider interface {
	GetWithToken([]byte) (UserProfile, error)
	Get(*http.Request) (UserProfile, error)
	Start()
	Stop()
}

type PasswordProvider interface {
	GetWithPassword(email string, password string) (UserProfile, error)
}

var providers = map[string]UserProfileProviderMaker{
	"worf": MakeWorfUserProfileProvider,
}

func MakeUserProfileProvider(settings map[string]any) (UserProfileProvider, error) {
	providerType, ok := settings["type"].(string)
	if !ok {
		return nil, fmt.Errorf("Provider type config missing (user-profile-provider.type)")
	}

	maker, ok := providers[providerType]

	if !ok {
		return nil, fmt.Errorf("Unknown provider type: %s", providerType)
	}

	return maker(settings)
}
