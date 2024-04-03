package main

import (
	"flag"
	"fmt"
	. "github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites"
	"github.com/klaro-org/sites/models"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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
		runServer()
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

	db, err := orm.Connect("linearize", settings.Database)

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

func runServer() {

	server := MakeServer(&App{
		StaticPrefix: "/static",
		Root: func(c Context) Element {

			return &HTMLElement{
				Tag:      "div",
				Children: []any{P("test")},
			}
		},
	})
	if err := server.Start(); err != nil {
		Log.Error("Cannot start server: %v", err)
	}
	wait()
}

func wait() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
}
