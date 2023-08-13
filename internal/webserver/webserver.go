package webserver

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type Webserver struct {
	tls    bool
	server *http.Server
}

func New(addr string, handler http.Handler) *Webserver {
	return &Webserver{
		server: &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,
			Addr:         addr,
			Handler:      handler,
		},
	}
}

func NewTLS(addr string, handler http.Handler, cacheDir string, hostNames []string) *Webserver {
	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(cacheDir),
		Email:      "dmytro.vovk@pm.me",
		HostPolicy: autocert.HostWhitelist(hostNames...),
	}

	config := manager.TLSConfig()
	config.MinVersion = tls.VersionTLS12
	config.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}
	config.CurvePreferences = []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256}

	return &Webserver{
		tls: true,
		server: &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,
			Addr:         addr,
			Handler:      manager.HTTPHandler(handler),
			TLSConfig:    config,
		},
	}
}

func (w *Webserver) Serve(name string) (err error) {
	if w.tls {
		log.Printf("%s starting at https://%s", name, w.server.Addr)
		err = w.server.ListenAndServeTLS("", "")
	} else {
		log.Printf("%s starting at http://%s", name, w.server.Addr)
		err = w.server.ListenAndServe()
	}

	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	return
}

func (w *Webserver) Stop(ctx context.Context) error {
	return w.server.Shutdown(ctx)
}
