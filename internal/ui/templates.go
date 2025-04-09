package ui

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippetbox/internal/modules/user/domain"
	"snippetbox/ui"
	"time"
)

type TemplateData struct {
	CurrentYear     int
	Model           any
	List            any
	Form            any
	Flash           string
	IsAuthenticated bool
	User            *domain.UserModel
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{"humanDate": humanDate}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"html/base.tmpl.html",
			"html/partials/nav.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
