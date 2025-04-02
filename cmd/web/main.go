package main

import (
	"crypto/tls"
	"flag"
	"net/http"
	"time"
)

func init() {
	Config.init()
	App.init()
}

func main() {
	flag.Parse()

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         Config.addr,
		ErrorLog:     App.errorLog,
		Handler:      App.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	App.infoLog.Println("Starting server on ", Config.addr)

	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	App.errorLog.Fatal(err)
}
