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

// WrapHTTPS wraps an http server in HTTPS TLS
func WrapHTTPS(s *http.Server, tlsAddr, cacheDir, host string) Server {
	httpsSrv := &HTTPSServer{
		Handler:     s.Handler,
		Addr:        tlsAddr,
		TLSDataDir:  cacheDir,
		AllowedHost: host,
	}
	m := httpsSrv.InitAutocert()
	httpSrv := &HTTPServer{
		Handler: makeHTTPToHTTPSRedirectMux(m),
		Addr:    s.Addr,
	}

	return &RedirectHTTPSServer{
		HTTPServer:  httpSrv,
		HTTPSServer: httpsSrv,
	}
}

// RedirectHTTPSServer redirects traffic going to the httpServer to the https server
type RedirectHTTPSServer struct {
	*HTTPSServer
	*HTTPServer
}

// ListenAndServe listen and serves on the calling server
func (s *RedirectHTTPSServer) ListenAndServe() error {
	errChan := make(chan error)
	go func() { errChan <- s.HTTPSServer.ListenAndServe() }()
	go func() { errChan <- s.HTTPServer.ListenAndServe() }()
	return <-errChan
}

// HTTPSServer represents everything needed to run an https server
type HTTPSServer struct {
	Handler     http.Handler
	Addr        string
	TLSDataDir  string
	AllowedHost string
	srv         *http.Server
}

// InitAutocert configures the http server and returns the autocert manager
func (s *HTTPSServer) InitAutocert() *autocert.Manager {
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
	m.GetCertificate(&tls.ClientHelloInfo{})
	s.srv = &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, s.Handler),
		Addr:         s.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
	}
	return m
}

// ListenAndServe listens and serves on the port of the calling server
func (s *HTTPSServer) ListenAndServe() error {
	return s.srv.ListenAndServeTLS("", "")
}

// HTTPServer represents everything needed to run an http server
type HTTPServer struct {
	Handler http.Handler
	Addr    string
}

// ListenAndServe listens and serves on the calling server's port
func (s *HTTPServer) ListenAndServe() error {
	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, s.Handler),
		Addr:         s.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}
