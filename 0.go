package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"time"

	"log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"fmt"
)

var counter int     					// Global Counter <- The Order
var tickcounter int 					// Global Counter <- AllTix
var upgrader = websocket.Upgrader{		// Websocket API v0.1
	ReadBufferSize:  64,				//  Default: 1024 ?
	WriteBufferSize: 64,				//  Default: 1024 ?
}

func main() {
	var entry string
	var static string
	var port string

	flag.StringVar(&entry, "entry", "../r3x/index.html", "R3X Entrypoint.")
	flag.StringVar(&static, "static", "../r3x/", "R3X Static Files Directory to serve.")
	flag.StringVar(&port, "port", "8000", "R3X Server Port.")
	flag.Parse()

	// <<< ROUTE >>>
	r := mux.NewRouter()
	// /api/v1/... production mindset etc...
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/ticker", TickerHandler).Methods("GET") // Ticker BTC - AUD
	api.HandleFunc("/buy", BuyHandler).Methods("POST")      // BTC Purchase
	// static files folder css/js...
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(static)))
	// 		*** not needed when static assets  ==  "/"
	// JavaScript entry-point (index.html).
	//r.PathPrefix("/index.html").HandlerFunc(IndexHandler(entry))

	// < < SERVE > >
	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Todo: 	Official Cert - ListenAndServeTLS(...)
	// https://gist.github.com/d-schmidt/587ceec34ce1334a5e60
	log.Fatal(srv.ListenAndServe())
}

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}

// UPDATE THE BTC - AUD PRICE
func BTC(){

}

// >> > >> > > > > > > >> >  >>  > ** * ** * * *
// 			BTC - AUD 		price 	ticker !  >  >    >
// >> > >> > > > > > > >> >  >>  > ** * ** * * *
func TickerHandler(w http.ResponseWriter, r *http.Request) {
	// Todo: ResponseHandler pls + performance analysis for each conn... w/ ping/pong or extended JSON timeVal diff based latency
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("--> Client Websocket Upgrade [fail]")
		log.Println("-->", err)
		return
	}

	// Check client WebSocket Todo: ( if no WS, REST API Fallback )
	if websocket.IsWebSocketUpgrade(r) {
		log.Println("---> Client Websocket Upgrade [OK]")
		log.Println("---> ", conn.LocalAddr().String())
	} // else REST


	data := map[string]interface{}{
		"id": "12345examplebro123",
		"ts": time.Now().Format(time.RFC3339),
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	conn.WriteJSON(b)
	//w.Write(b)
}

// ===================================
// ===================================
// ! ! ! 	BITCOIN ORDER 		 ! ! !
// ===================================
func BuyHandler(w http.ResponseWriter, r *http.Request) {
	amount := r.PostForm.Get("btc") // Amount requested !

	data := map[string]interface{}{
		"id":  fmt.Sprintf("%d", counter),
		"btc": amount,
		"ts":  time.Now().Format(time.RFC3339),
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	counter++ // The Order Counter
	w.Write(b)
}
