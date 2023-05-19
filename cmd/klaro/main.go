package main

import (
	"fmt"
	"github.com/klaro-org/klaro-cms"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	backend := klaro.MakeInMemoryBackend()

	server := klaro.MakeServer(klaro.Options{
		StaticPrefix: "/static",
	}, backend)
	server.Start()
	wait()
}

func wait() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
}
