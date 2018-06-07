package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jgall/easyhttps"
)

var (
	flgHTTPS        = true
	directory       = "test"
	host            = ""
	httpsPort       = "443"
	httpPort        = "80"
	tlsCertCacheDir = "/.tlsCertCache"
)

func parseFlags() {
	flag.BoolVar(&flgHTTPS, "https", true, "if true, we start HTTPS server")
	flag.StringVar(&directory, "dir", "test", "the directory from which to serve files")
	flag.StringVar(&tlsCertCacheDir, "cache-dir", "/.tlsCertCache", "a cache directory for ")
	flag.StringVar(&host, "host", "", "the host you're running on")
	flag.StringVar(&httpsPort, "httpsPort", "443", "https listening port")
	flag.StringVar(&httpPort, "httpPort", "80", "http listening port")
	flag.Parse()
}

func main() {
	parseFlags()
	srv := &http.Server{
		Handler: http.FileServer(http.Dir(directory)),
		Addr:    ":" + httpPort,
	}

	if flgHTTPS {
		httpsSrv := easyhttps.WrapHTTPS(srv, ":"+httpsPort, tlsCertCacheDir, host)
		log.Fatal(httpsSrv.ListenAndServe())
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}
