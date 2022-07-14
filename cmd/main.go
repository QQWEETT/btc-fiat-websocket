package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"quest/internals/app"
	"quest/internals/cfg"
)

func main() {
	config := cfg.LoadAndStoreConfig()
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server := app.NewServer(config, ctx)
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		server.Shutdown()
		cancel()
	}()
	server.Serve()
}
