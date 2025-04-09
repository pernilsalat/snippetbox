package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"net/http"
	"snippetbox/internal/modules/user/repository"
	"snippetbox/internal/platform"
	"snippetbox/internal/platform/database"
)

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-security-policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("x-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func LogRequest(app *Application) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

			next.ServeHTTP(w, r)
		})
	}
}

func RecoverPanic(app *Application) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					app.ServerError(w, fmt.Errorf("%s", err))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthentication(app *Application) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !app.IsAuthenticated(r) {
				http.Redirect(w, r, "/user/login", http.StatusSeeOther)
				return
			}
			w.Header().Set("Cache-Control", "no-store")

			next.ServeHTTP(w, r)
		})
	}
}

func Authenticate(app *Application) alice.Constructor {
	ur := repository.NewUser(app.DB)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			id := app.SessionManager.GetInt(r.Context(), "authenticatedUserID")
			if id == 0 {
				next.ServeHTTP(w, r)
				return
			}
			//exists, err := ur.Exists(id)
			user, err := ur.Get(id)
			if !errors.Is(err, database.ErrNoRecord) && err != nil {
				app.ServerError(w, err)
				return
			} else if user != nil {
				ctx := context.WithValue(r.Context(), platform.IsAuthenticatedContextKey, true)
				ctx = context.WithValue(ctx, platform.UserContextKey, user)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func Dynamic(app *Application) alice.Chain {
	return alice.New(app.SessionManager.LoadAndSave, NoSurf, Authenticate(app))
}

func Protected(app *Application) alice.Chain {
	return Dynamic(app).Append(RequireAuthentication(app))
}

func Standard(app *Application) alice.Chain {
	return alice.New(RecoverPanic(app), LogRequest(app), SecureHeaders)
}
