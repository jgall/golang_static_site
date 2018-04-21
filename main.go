package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
)

func main() {
	port := flag.String("port", "8000", "the `port` to listen on.")
	entry := flag.String("entry", "entry", "the directory to serve static files from.")
	flag.Parse()

	r := http.NewServeMux()
	fs := http.FileServer(http.Dir(*entry))
	r.Handle("/", fs)
	log.Printf("Serving files from %v on port %v...\n", *entry, *port)

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         "0.0.0.0:" + *port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
