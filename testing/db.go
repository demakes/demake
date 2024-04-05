package testing

import (
	"errors"
	"fmt"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites"
	"github.com/klaro-org/sites/models"
	"os"
	"regexp"
)

var urlRegexp = regexp.MustCompile(`^(.*)\?`)

func DB(settings *sites.Settings) (orm.DB, error) {

	if !settings.Test {
		return nil, fmt.Errorf("not in test mode, aborting DB setup!")
	}

	if settings.Database.Type == "sqlite3" {
		// we remove the SQLite database
		if path := urlRegexp.FindStringSubmatch(settings.Database.Url); path == nil {
			return nil, fmt.Errorf("cannot parse path")
		} else {
			if _, err := os.Stat(path[1]); err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return nil, fmt.Errorf("cannot stat file: %w", err)
				}
			} else if err := os.Remove(path[1]); err != nil {
				return nil, fmt.Errorf("cannot remove database")
			}
		}
	}

	// we clear all previous connections
	if err := orm.ClearConnections(); err != nil {
		return nil, err
	}

	db, err := orm.Connect("klaro", settings.Database)

	if err != nil {
		return nil, err
	}

	manager, err := orm.MakeMigrationManager("migrations", models.Migrations, db, settings.Database.Type)

	if err != nil {
		return nil, fmt.Errorf("cannot create migration manager: %w", err)
	}

	return db, manager.Migrate(manager.LatestVersion())
}
