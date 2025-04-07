package main

import (
	"flag"
	"net/http"
	"snippetbox/internal/platform"
	"snippetbox/internal/platform/database/mysql"
	"snippetbox/internal/platform/web"
	"time"
)

func init() {
	platform.Config.Init()
	//platform.App.Init()
}

func main() {
	flag.Parse()

	db, err := mysql.Connect(platform.Config.Dsn)
	if err != nil {
		platform.ErrorLog.Fatal(err)
	}

	app, err := web.NewApplication(db)
	if err != nil {
		platform.ErrorLog.Fatal(err)
	}

	srv := &http.Server{
		Addr:         platform.Config.Addr,
		ErrorLog:     app.ErrorLog,
		Handler:      platform.App.Routes(),
		TLSConfig:    platform.Config.TlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.InfoLog.Println("Starting server on ", platform.Config.Addr)

	err = srv.ListenAndServeTLS(
		platform.Config.TlsDir+"/cert.pem",
		platform.Config.TlsDir+"/key.pem",
	)

	app.ErrorLog.Fatal(err)
}
