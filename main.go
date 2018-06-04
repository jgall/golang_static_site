package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jgall/golang_static_site/pkg/server"
)

var (
	flgHTTPS  = true
	directory = "test"
	host      = ""
	httpsPort = "443"
	httpPort  = "80"
)

func makeHTTPToHTTPSRedirectMux() *http.ServeMux {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return mux
}

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

	if flgHTTPS {
		httpsSrv := &server.HTTPSServer{
			Mux:         makeMainMux(directory),
			Port:        httpsPort,
			TLSDataDir:  ".",
			AllowedHost: host,
		}
		httpSrv := &server.HTTPServer{
			Mux:  makeHTTPToHTTPSRedirectMux(),
			Port: httpPort,
		}
		errChan := make(chan error)
		go func() { errChan <- httpsSrv.ListenAndServe() }()
		go func() { errChan <- httpSrv.ListenAndServe() }()
		if err := <-errChan; err != nil {
			// Quit on the first error that comes through the error channel
			log.Fatal(err)
			panic(err)
		}
	} else {
		httpSrv := server.HTTPServer{
			Mux:  makeMainMux(directory),
			Port: httpPort,
		}
		log.Fatal(httpSrv.ListenAndServe())
	}
}
