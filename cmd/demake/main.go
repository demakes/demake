package main

import (
	"flag"
	"fmt"
	"github.com/gospel-sh/gospel/orm"
	"github.com/demakes/demake"
	"github.com/demakes/demake/models"
	"log/slog"
	"os"
)

func main() {

	migrateFlags := flag.NewFlagSet("migrate", flag.ExitOnError)

	var cmd string

	if len(os.Args) < 2 {
		cmd = "run"
	} else {
		cmd = os.Args[1]
	}

	switch cmd {
	case "migrate":
		migrateFlags.Parse(os.Args[2:])
		if err := runMigrations(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(-1)
		}
	case "run":
		if err := sites.Run(); err != nil {
			fmt.Printf("error running: %v", err)
			os.Exit(-1)
		}
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		os.Exit(-1)
	}
}

func runMigrations() error {

	settings, err := sites.LoadSettings()

	if err != nil {
		return err
	}

	db, err := orm.Connect("klaro", settings.Database)

	if err != nil {
		return err
	}

	manager, err := orm.MakeMigrationManager("migrations", models.Migrations, db, settings.Database.Type)

	if err != nil {
		return err
	}

	slog.Info("Running migrations...", slog.Int("version", manager.LatestVersion()))
	return manager.Migrate(manager.LatestVersion())

}
