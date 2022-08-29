package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"golang.org/x/crypto/acme/autocert"
)

const (
	idleTimeout  = 60 * time.Second
	readTimeout  = 60 * time.Second
	writeTimeout = 60 * time.Second
	cacheDir     = "certs"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {

	var TLSConfig *tls.Config

	if cfg.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache(cacheDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(strings.Join(cfg.HostWhitelist, ",")),
		}

		TLSConfig = manager.TLSConfig()
	}

	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         cfg.ServerAddr,
			Handler:      handler,
			TLSConfig:    TLSConfig,
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

func (s *Server) Run() error {
	if s.cfg.EnableHTTPS {
		return s.httpServer.ListenAndServeTLS("", "")
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
