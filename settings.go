package sites

import (
	"encoding/json"
	"github.com/gospel-sh/gospel/orm"
	"os"
)

type Settings struct {
	Test     bool                  `json:"test"`
	Database *orm.DatabaseSettings `json:"database"`
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
