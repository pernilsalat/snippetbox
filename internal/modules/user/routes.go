package user

import (
	"net/http"
	"snippetbox/internal/modules/user/handler"
	"snippetbox/internal/modules/user/repository"
	"snippetbox/internal/modules/user/service"
	"snippetbox/internal/platform/web"
)

func InitRoutes(app *web.Application) {
	ur := repository.NewUser(app.DB)
	us := service.NewUser(ur)
	uh := handler.NewUser(app, us)

	dynamic := web.Dynamic(app)
	app.Router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(uh.UserSignup))
	app.Router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(uh.UserSignupPost))
	app.Router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(uh.UserLogin))
	app.Router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(uh.UserLoginPost))

	protected := web.Protected(app)
	app.Router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(uh.UserLogoutPost))
	app.Router.Handler(http.MethodGet, "/user/view", protected.ThenFunc(uh.UserView))
}
