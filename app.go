package sites

import (
	"fmt"
	. "github.com/gospel-sh/gospel"
	"github.com/gospel-sh/gospel/orm"
	"github.com/demakes/demake/models"
	"github.com/demakes/demake/ui"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type MainServer struct {
	settings  *Settings
	appServer *Server
	db        orm.DB
}

func (m *MainServer) ServeSite(site *models.Site, w http.ResponseWriter, r *http.Request) {
	appServer := MakeServer(&App{
		Root:         ui.ServeSite(m.db, site),
		StaticPrefix: "/static",
	})

	appServer.ServeHTTP(w, r)

}

func (m *MainServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dbf := func() orm.DB { return m.db }

	if strings.HasPrefix(r.URL.Path, "/demake") {
		// we serve the admin UI
		m.appServer.ServeHTTP(w, r)
		return
	}

	site := &models.Site{}

	orm.Init(site, dbf)

	hostnameAndPort := r.Host
	hostname := strings.Split(hostnameAndPort, ":")[0]

	if err := site.ByHostname(hostname); err != nil {
		if err == orm.NotFound {
			fmt.Fprintf(w, "unknown site")
			w.Header().Add("content-type", "text/plain")
			w.WriteHeader(404)
			return
		}
		// this is an unexpected error
		fmt.Fprintf(w, "error loading site")
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(500)
		return
	}

	m.ServeSite(site, w, r)
}

func Run() error {

	settings, err := LoadSettings()

	if err != nil {
		return err
	}

	profileProvider, err := MakeUserProfileProvider(settings.Auth)

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

	mainServer := &http.Server{
		Addr: ":8001",
		Handler: &MainServer{
			db:       db,
			settings: settings,
			appServer: MakeServer(&App{
				Root:         ui.Root(db, profileProvider),
				StaticPrefix: "/static",
			}),
		},
	}

	go mainServer.ListenAndServe()

	wait()

	return nil
}

func wait() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
}
