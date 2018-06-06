package easyhttps

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

func makeHTTPToHTTPSRedirectMux(m *autocert.Manager) http.HandlerFunc {
	hasHTTPS := false
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	callbackHandler := m.HTTPHandler(nil)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/.well-known/acme-challenge/") {
			callbackHandler.ServeHTTP(w, r)
			hasHTTPS = true
		} else if hasHTTPS {
			handleRedirect(w, r)
		} else {
			fmt.Fprint(w, "Waiting on TLS Certificate")
		}
	}
}
