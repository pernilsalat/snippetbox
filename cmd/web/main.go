package main

import (
	"flag"
	"net/http"
)

func init() {
	Config.init()
	App.init()
}

func main() {
	flag.Parse()

	srv := &http.Server{
		Addr:     Config.addr,
		ErrorLog: App.errorLog,
		Handler:  App.routes(),
	}

	App.infoLog.Println("Starting server on ", Config.addr)

	err := srv.ListenAndServe()

	App.errorLog.Fatal(err)
}
