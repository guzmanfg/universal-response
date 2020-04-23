package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

func StartServer(mux *http.ServeMux, addr string, sslAddr string, ssl map[string]string) chan error {
	if sslAddr == "" {
		return Http(mux, addr)
	}

	return Https(mux, addr, sslAddr, ssl)
}

func Http(mux *http.ServeMux, addr string) chan error {
	errs := make(chan error)
	if err := http.ListenAndServe(addr, mux); err != nil {
		errs <- err
	}
	return errs
}

func Https(mux *http.ServeMux, addr string, sslAddr string, ssl map[string]string) chan error {
	tlsConfig := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := &http.Server{
		Addr:         sslAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}

	return Run(srv, addr, sslAddr, ssl)
}

func Run(srv *http.Server, addr string, sslAddr string, ssl map[string]string) chan error {

	errs := make(chan error)

	host, _, err := net.SplitHostPort(sslAddr)
	if err == nil && host == "" {
		host = "https://localhost" + sslAddr
	} else {
		host = "https://" + sslAddr
	}

	// Starting HTTPS server
	go func() {
		if err := srv.ListenAndServeTLS(ssl["cert"], ssl["key"]); err != nil {
			errs <- err
		}
	}()

	// Starting HTTP server
	go func() {
		if err := http.ListenAndServe(addr, http.HandlerFunc(RedirectHttps(host))); err != nil {
			errs <- err
		}
	}()


	return errs
}

func RedirectHttps(sslAddr string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		target := sslAddr + req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		log.Printf("redirect to: %s", target)
		http.Redirect(w, req, target, http.StatusPermanentRedirect)
	}
}
