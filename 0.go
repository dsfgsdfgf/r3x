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
	"fmt"
)

func main() {
	var entry string
	var static string
	var port string

	flag.StringVar(&entry, "entry", "../r3x/index.html", "R3X Entrypoint.")
	flag.StringVar(&static, "static", "../r3x/", "R3X Static Files Directory to serve.")
	flag.StringVar(&port, "port", "8000", "R3X Server Port.")
	flag.Parse()

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/ticker", TickerHandler).Methods("GET")		// Ticker BTC - AUD
	api.HandleFunc("/buy", BuyHandler).Methods("POST")			// BTC Purchase

	// static files folder css/js...
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(static)))
	// 		*** not needed when static assets  ==  "/"
	// JavaScript entry-point (index.html).
	//r.PathPrefix("/index.html").HandlerFunc(IndexHandler(entry))

	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    "127.0.0.1:" + port,
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

// >> > >> > > > > > > >> >  >>  > ** * ** * * *
// 			BTC - AUD 		price 	ticker !  >  >    >
// >> > >> > > > > > > >> >  >>  > ** * ** * * *
func TickerHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"id": "12345",
		"ts": time.Now().Format(time.RFC3339),
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write(b)
}


// ===================================
// ===================================
// ! ! ! 	BITCOIN ORDER 		 ! ! !
// ===================================
func BuyHandler(w http.ResponseWriter, r *http.Request) {
	var counter int;
	amount := r.PostForm.Get("btc")		// Amount requested !

	data := map[string]interface{}{
		"id": fmt.Printf("%d", counter),
		"btc": amount,
		"ts": time.Now().Format(time.RFC3339),
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write(b)
}
