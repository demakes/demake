package sites

import (
	"fmt"
	. "github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites/auth"
	"github.com/klaro-org/sites/ui"
	"os"
	"os/signal"
	"syscall"
)

var userProfileSettings = map[string]any{
	"type": "worf",
	"url":  "http://localhost:5000/v1",
}

func Run() error {

	profileProvider, err := auth.MakeUserProfileProvider(userProfileSettings)

	if err != nil {
		return err
	}

	settings, err := LoadSettings()

	if err != nil {
		return err
	}

	db, err := orm.Connect("klaro", settings.Database)

	if err != nil {
		return err
	}

	if settings.Database.Type == "sqlite3" {
		// this enables WAL mode, which drastically speeds up execution
		if _, err := db.Exec(`PRAGMA journal_mode = WAL; PRAGMA synchronous = NORMAL;`); err != nil {
			return err
		}
	}

	server := MakeServer(&App{
		Root:         ui.Root(db, profileProvider),
		StaticPrefix: "/static",
	})

	if err := server.Start(); err != nil {
		return err
	}

	wait()

	return nil
}

func wait() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
}
