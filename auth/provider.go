package auth

import (
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
