package sites

import (
	"encoding/json"
	"fmt"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites/auth"
	"os"
)

type Settings struct {
	Test     bool                  `json:"test"`
	Database *orm.DatabaseSettings `json:"database"`
	Auth     *AuthSettings         `json:"auth"`
}

type AuthSettings struct {
	Type   string               `json:"type"`
	Worf   *auth.WorfSettings   `json:"worf"`
	Simple *auth.SimpleSettings `json:"simple"`
}

func LoadSettings() (*Settings, error) {

	settingsPath := os.Getenv("KLARO_SETTINGS")

	if settingsPath == "" {
		settingsPath = "settings/dev/sqlite.json"
	}

	if config, err := os.ReadFile(settingsPath); err != nil {
		return nil, err
	} else {
		var settings *Settings
		if err := json.Unmarshal(config, &settings); err != nil {
			return nil, err
		}
		return settings, nil
	}
}

func MakeUserProfileProvider(settings *AuthSettings) (auth.UserProfileProvider, error) {
	switch settings.Type {
	case "worf":
		return auth.MakeWorfUserProfileProvider(settings.Worf)
	case "simple":
		return auth.MakeSimpleUserProfileProvider(settings.Simple)
	}

	return nil, fmt.Errorf("unknown user profile provider: %s", settings.Type)
}
