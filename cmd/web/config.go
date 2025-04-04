package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"os"
	"snippetbox/internal/models"
	"time"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

func (c *config) init() {
	flag.StringVar(&c.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&c.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&c.dsn, "dsn", "user:password@/db?parseTime=true", "MySQL data source name")
}

var Config = &config{}

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.ISnippetModel
	users          models.IUserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func (app *application) init() {
	var (
		db  *sql.DB
		err error
	)
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if db, err = openDb(Config.dsn); err != nil {
		app.errorLog.Fatal(err)
	}

	app.snippets = &models.SnippetModel{DB: db}
	app.users = &models.UserModel{DB: db}

	if app.templateCache, err = newTemplateCache(); err != nil {
		app.errorLog.Fatal(err)
	}
	app.formDecoder = form.NewDecoder()

	app.sessionManager = scs.New()
	app.sessionManager.Store = mysqlstore.New(db)
	app.sessionManager.Lifetime = 12 * time.Hour
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

var App = &application{}
