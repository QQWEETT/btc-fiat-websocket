package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"quest/internals/app/handlers"
)

func CreateRoutes(btcHandler *handlers.BTCHandler, fiatHandler *handlers.FiatHandler, btcFiatHandler *handlers.BTCFiatHandler, websocket *handlers.WebsocketHandler) *mux.Router {
	r := mux.NewRouter()
	pool := websocket.NewPool()
	go pool.Start()
	r.HandleFunc("/api/btcusdt", btcHandler.LastBTC).Methods("GET")
	r.HandleFunc("/api/btcusdt", btcHandler.FindBTC).Methods("POST")
	r.HandleFunc("/api/currencies", fiatHandler.LastFiat).Methods("GET")
	r.HandleFunc("/api/currencies", fiatHandler.FindFiat).Methods("POST")
	r.HandleFunc("/api/latest", btcFiatHandler.SocketHandler).Methods("GET")
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(pool, w, r)
	})
	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFound).GetHandler()
	go btcHandler.CreateBTC()
	go fiatHandler.CreateFiat()
	go btcFiatHandler.CreateBTCFiat()
	return r
}
