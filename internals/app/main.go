package app

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"quest/api"
	"quest/api/middleware"
	db3 "quest/internals/app/db"
	"quest/internals/app/handlers"
	"quest/internals/app/processors"
	"quest/internals/cfg"
	"time"
)

type Server struct {
	config cfg.Cfg
	ctx    context.Context
	srv    *http.Server
	db     *pgxpool.Pool
}

func NewServer(config cfg.Cfg, ctx context.Context) *Server {
	server := new(Server)
	server.ctx = ctx
	server.config = config
	return server

}

func (server *Server) Serve() {
	log.Println("Starting server")
	var err error
	server.db, err = pgxpool.Connect(server.ctx, server.config.GetDBString())
	if err != nil {
		log.Fatalln(err)
	}
	btcStorage := db3.NewBTCStorage(server.db)
	fiatStorage := db3.NewFiatStorage(server.db)
	BTCFiatStorage := db3.NewBTCFiatStorage(server.db)
	WebsocketStorage := db3.NewWebsocketStorage(server.db)

	btcProcessor := processors.NewBTCProcessor(btcStorage)
	fiatProcessor := processors.NewFiatProcessor(fiatStorage)
	BTCFiatProcessor := processors.NewBTCFiatProcessor(BTCFiatStorage)
	WebsocketProcessor := processors.NewWebsocketProcessor(WebsocketStorage)

	btcHandler := handlers.NewBTCHandler(btcProcessor)
	fiatHandler := handlers.NewFiatHandler(fiatProcessor)
	BTCFiatHandler := handlers.NewBTCFiatHandler(BTCFiatProcessor)
	WebsocketHandler := handlers.NewWebsocketHandler(WebsocketProcessor)

	routes := api.CreateRoutes(btcHandler, fiatHandler, BTCFiatHandler, WebsocketHandler)
	routes.Use(middleware.RequestLog)

	server.srv = &http.Server{
		Addr:    ":" + server.config.Port,
		Handler: routes,
	}

	log.Println("Server started")

	err = server.srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

	return
}

func (server *Server) Shutdown() {
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.db.Close()
	defer func() {
		cancel()
	}()
	var err error
	if err = server.srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown failed:%e ", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
