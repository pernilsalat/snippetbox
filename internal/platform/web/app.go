package web

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"snippetbox/internal/modules/user/domain"
	"snippetbox/internal/platform"
	"snippetbox/internal/platform/database"
	"snippetbox/internal/ui"
	"time"
)

type Application struct {
	Router         *httprouter.Router
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
	DB             database.DB
}

func NewApplication(db *sql.DB) (*Application, error) {
	templateCache, err := ui.NewTemplateCache()
	if err != nil {
		return nil, err
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	return &Application{
		Router:         httprouter.New(),
		InfoLog:        platform.InfoLog,
		ErrorLog:       platform.ErrorLog,
		TemplateCache:  templateCache,
		FormDecoder:    form.NewDecoder(),
		SessionManager: sessionManager,
		DB:             db,
	}, nil
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.ErrorLog.Output(2, trace)

	if platform.Config.Debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) Render(w http.ResponseWriter, status int, page string, td *ui.TemplateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(w, err)
		return
	}
	buffer := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buffer, "base", td); err != nil {
		app.ServerError(w, err)
		return
	}

	w.WriteHeader(status)
	if _, err := buffer.WriteTo(w); err != nil {
		app.ServerError(w, err)
	}
}

func (app *Application) NewTemplateData(r *http.Request) *ui.TemplateData {
	return &ui.TemplateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.IsAuthenticated(r),
		UserName:        app.GetCurrentUser(r).Name,
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *Application) DecodePostForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := app.FormDecoder.Decode(dst, r.PostForm); err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			// TODO: use ServerError instead
			panic(err)
		}
		return err
	}

	return nil

}

func (app *Application) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(platform.IsAuthenticatedContextKey).(bool)

	return ok && isAuthenticated
}

func (app *Application) GetCurrentUser(r *http.Request) *domain.UserModel {
	user, ok := r.Context().Value(platform.UserContextKey).(*domain.UserModel)

	if !ok {
		return &domain.UserModel{}
	}
	return user
}
