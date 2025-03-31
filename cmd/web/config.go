package main

import (
	"flag"
	"log"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

func (c *config) init() {
	flag.StringVar(&c.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&c.staticDir, "static-dir", "./ui/static", "Path to static assets")
}

var Config = &config{}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func (app *application) init() {
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
}

var App = &application{}
