package easyhttps

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

// Server represents the things a server can do
type Server interface {
	ListenAndServe() error
}

// WrapHTTPS wraps an http server in HTTPS TLS
func WrapHTTPS(s *http.Server, tlsAddr, cacheDir, host string) Server {
	httpsSrv := &HTTPSServer{
		TLSDataDir:  cacheDir,
		AllowedHost: host,
		srv: &http.Server{
			Handler:      s.Handler,
			Addr:         tlsAddr,
			WriteTimeout: s.WriteTimeout,
			ReadTimeout:  s.ReadTimeout,
			IdleTimeout:  s.IdleTimeout,
		},
	}
	m := httpsSrv.InitAutocert()
	httpSrv := &HTTPServer{
		srv: &http.Server{
			Handler:      httpsSrv.makeHTTPToHTTPSRedirectMux(m),
			Addr:         s.Addr,
			WriteTimeout: s.WriteTimeout,
			ReadTimeout:  s.ReadTimeout,
			IdleTimeout:  s.IdleTimeout,
		},
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
	if err := os.MkdirAll(s.TLSDataDir, 0700); err != nil {
		log.Printf("warning: autocert.NewListener not using a cache: %v", err)
		return nil
	}
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(s.TLSDataDir),
	}
	m.GetCertificate(&tls.ClientHelloInfo{})
	s.srv.TLSConfig = &tls.Config{
		GetCertificate: m.GetCertificate,
	}

	return m
}

// ListenAndServe listens and serves on the port of the calling server
func (s *HTTPSServer) ListenAndServe() error {
	return s.srv.ListenAndServeTLS("", "")
}

// HTTPServer represents everything needed to run an http server
type HTTPServer struct {
	srv *http.Server
}

// ListenAndServe listens and serves on the calling server's port
func (s *HTTPServer) ListenAndServe() error {
	return s.srv.ListenAndServe()
}
