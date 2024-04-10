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

func runServer() {

	el := Div(
		P("This is a test"),
		func(c Context) (any, error) {
			return P("another one"), nil
		},
		Strong("strong"),
		Route("/test(/[a-z]+)?", Strong("another test")),
	)

	server := MakeServer(&App{
		StaticPrefix: "/static",
		Root: func(c Context) Element {
			if v, err := el.Generate(c); err != nil {
				return Div(Fmt("error: %v", err))
			} else if el, ok := v.(Element); !ok {
				return Div("not an element")
			} else {
				return el
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