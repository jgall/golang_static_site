package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"golang.org/x/crypto/acme/autocert"
)

// Server represents the things a server can do
type Server interface {
	ListenAndServe() error
}

// HTTPSServer represents everything needed to run an https server
type HTTPSServer struct {
	Mux         *http.ServeMux
	Port        string
	TLSDataDir  string
	AllowedHost string
}

// ListenAndServe listens and serves on the port of the calling server
func (s *HTTPSServer) ListenAndServe() error {
	hostPolicy := func(ctx context.Context, host string) error {
		if host == s.AllowedHost {
			return nil
		}
		return fmt.Errorf("acme/autocert: only %s host is allowed", s.AllowedHost)
	}
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(s.TLSDataDir),
	}
	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, s.Mux),
		Addr:         ":" + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
	}
	return srv.ListenAndServeTLS("", "")
}

// HTTPServer represents everything needed to run an http server
type HTTPServer struct {
	Mux  *http.ServeMux
	Port string
}

// ListenAndServe listens and serves on the calling server's port
func (s *HTTPServer) ListenAndServe() error {
	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, s.Mux),
		Addr:         ":" + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}
