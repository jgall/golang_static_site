package main

// https://blog.kowalczyk.info/article/Jl3G/https-for-free-in-go.html
// To run:
// go run main.go
// Command-line options:
//   -production : enables HTTPS on port 443
//   -redirect-to-https : redirect HTTP to HTTTPS

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
	flag.Parse()
}

func main() {
	parseFlags()

	if flgHTTPS {
		httpsSrv := &server.HTTPSServer{
			Mux:         makeMainMux(),
			Port:        "443",
			TLSDataDir:  ".",
			AllowedHost: host,
		}
		httpSrv := &server.HTTPServer{
			Mux:  makeHTTPToHTTPSRedirectMux(),
			Port: "80",
		}
		go log.Fatal(httpsSrv.ListenAndServe())
		go log.Fatal(httpSrv.ListenAndServe())
	} else {
		httpSrv := server.HTTPServer{
			Mux:  makeMainMux(),
			Port: "80",
		}
		log.Fatal(httpSrv.ListenAndServe())
	}
}
