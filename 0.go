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
)

func main() {
	var entry string
	var static string
	var port string

	flag.StringVar(&entry, "entry", "./index.html", "the entrypoint to serve.")
	flag.StringVar(&static, "static", ".", "the directory to serve static files from.")
	flag.StringVar(&port, "port", "8000", "the `port` to listen on.")
	flag.Parse()

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/users", GetUsersHandler).Methods("GET")
	// Use a custom 404 handler for API paths.
	api.NotFoundHandler = JSONNotFound

	// static assets
	r.PathPrefix("/dist").Handler(http.FileServer(http.Dir(static)))
	//  JavaScript entry-point (index.html).
	r.PathPrefix("/").HandlerFunc(IndexHandler(entry))

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

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
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
