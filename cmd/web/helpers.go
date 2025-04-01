package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, td *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	buffer := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buffer, "base", td); err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	if _, err := buffer.WriteTo(w); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) newTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}

}
