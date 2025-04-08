package modules

import (
	"net/http"
	"snippetbox/internal/modules/ping"
	"snippetbox/internal/modules/snippet"
	"snippetbox/internal/modules/user"
	"snippetbox/internal/platform/web"
	"snippetbox/ui"
)

func Init(app *web.Application) {
	app.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	app.Router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	// NOTE: Simple handler for testing
	app.Router.HandlerFunc(http.MethodGet, "/ping", ping.Handler)

	user.InitRoutes(app)
	snippet.InitRoutes(app)
}
