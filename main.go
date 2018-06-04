package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/jgall/golang_static_site/pkg/server"
)

const (
	htmlIndex = `<html><body>Welcome fren</body></html>`
)

var (
	flgHTTPS  = true
	directory = "test"
	host      = ""
	httpsPort = "443"
	httpPort  = "80"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlIndex)
}

func makeHTTPToHTTPSRedirectMux() *http.ServeMux {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return mux
}

func makeMainMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
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
			Mux:         makeMainMux(),
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
			Mux:  makeMainMux(),
			Port: httpPort,
		}
		log.Fatal(httpSrv.ListenAndServe())
	}
}
