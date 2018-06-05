package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jgall/easyhttps"
)

var (
	flgHTTPS  = true
	directory = "test"
	host      = ""
	httpsPort = "443"
	httpPort  = "80"
)

func makeMainMux(dir string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(dir)))
	return mux
}

func parseFlags() {
	flag.BoolVar(&flgHTTPS, "https", true, "if true, we start HTTPS server")
	flag.StringVar(&directory, "dir", "test", "the directory from which to serve files")
	flag.StringVar(&host, "host", "", "the host you're running on")
	flag.StringVar(&httpsPort, "httpsPort", "443", "https listening port")
	flag.StringVar(&httpPort, "httpPort", "80", "http listening port")
	flag.Parse()
}

func main() {
	parseFlags()
	srv := &http.Server{
		Handler: makeMainMux(directory),
		Addr:    ":" + httpPort,
	}

	if flgHTTPS {
		httpsSrv := easyhttps.WrapHTTPS(srv, ":"+httpsPort, ".", host)
		log.Fatal(httpsSrv.ListenAndServe())
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}
