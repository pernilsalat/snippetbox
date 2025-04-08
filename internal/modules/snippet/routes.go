package snippet

import (
	"net/http"
	"snippetbox/internal/modules/snippet/handler"
	"snippetbox/internal/modules/snippet/repository"
	"snippetbox/internal/modules/snippet/service"
	"snippetbox/internal/platform/web"
)

func InitRoutes(app *web.Application) {
	// Initialize the Snippet module
	// This function can be used to set up routes, handlers, etc.
	sr := repository.NewSnippet(app.DB)
	ss := service.NewSnippet(sr)
	sh := handler.NewSnippet(app, ss)

	dynamic := web.Dynamic(app)
	app.Router.Handler(http.MethodGet, "/", dynamic.ThenFunc(sh.SnipperList))
	app.Router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(sh.SnippetView))

	protected := web.Protected(app)
	app.Router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(sh.SnippetCreateView))
	app.Router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(sh.SnippetCreatePost))
}
