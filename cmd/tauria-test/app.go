package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nichojovi/tauria-test/cmd/config"
	"github.com/nichojovi/tauria-test/cmd/internal"
	"github.com/nichojovi/tauria-test/cmd/web"
	"github.com/nichojovi/tauria-test/internal/utils/auth"
	"github.com/nichojovi/tauria-test/internal/utils/database"
)

func main() {
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg, "main")

	//DATABASE
	datastore := database.New(*cfg, "mysql")

	//SERVICE
	service := internal.GetService(datastore, cfg)

	//AUTH
	auth := auth.New(&auth.Opts{UserService: service.User})

	server := web.New(&web.Opts{
		ListenAddress: cfg.Server.Port,
		Service:       service,
		AuthService:   auth,
	})
	go server.Run()

	select {
	case <-terminateSignal():
		log.Println("Exiting gracefully...")
	case err := <-server.ListenError():
		log.Println("Error starting web server, exiting gracefully:", err)
	}
}

func terminateSignal() chan os.Signal {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	return term
}
