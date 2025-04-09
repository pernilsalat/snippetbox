package handler

import (
	"errors"
	"net/http"
	"snippetbox/internal/modules/user/domain"
	"snippetbox/internal/modules/user/repository"
	"snippetbox/internal/modules/user/service"
	"snippetbox/internal/platform/web"
	"snippetbox/internal/validator"
)

type UserHandler struct {
	service *service.User
	app     *web.Application
}

func NewUser(app *web.Application, svc *service.User) *UserHandler {
	return &UserHandler{
		service: svc,
		app:     app,
	}
}

func (h *UserHandler) UserSignup(w http.ResponseWriter, r *http.Request) {
	td := h.app.NewTemplateData(r)
	td.Form = &domain.UsersSignupForm{}
	h.app.Render(w, http.StatusOK, "signup.tmpl.html", td)
}

func (h *UserHandler) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	var form domain.UsersSignupForm
	if err := h.app.DecodePostForm(r, &form); err != nil {
		h.app.ClientError(w, http.StatusBadRequest)
		return
	}
	err := h.service.UserSignup(&form)
	if errors.Is(err, validator.ErrInvalidForm) || errors.Is(err, repository.ErrDuplicateEmail) {
		td := h.app.NewTemplateData(r)
		td.Form = form

		h.app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", td)
		return
	} else if err != nil {
		h.app.ServerError(w, err)
		return
	}

	h.app.SessionManager.Put(r.Context(), "flash", "User created successfully")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	td := h.app.NewTemplateData(r)
	td.Form = &domain.UserLoginForm{}
	h.app.Render(w, http.StatusOK, "login.tmpl.html", td)
}

func (h *UserHandler) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	var form domain.UserLoginForm
	if err := h.app.DecodePostForm(r, &form); err != nil {
		h.app.ClientError(w, http.StatusBadRequest)
		return
	}

	id, err := h.service.Authenticate(&form)
	if errors.Is(err, validator.ErrInvalidForm) || errors.Is(err, repository.ErrInvalidCredentials) {
		td := h.app.NewTemplateData(r)
		td.Form = form
		h.app.Render(w, http.StatusUnprocessableEntity, "login.tmpl.html", td)
		return
	} else if err != nil {
		h.app.ServerError(w, err)
		return
	}

	err = h.app.SessionManager.RenewToken(r.Context())
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	redirectPath := h.app.SessionManager.PopString(r.Context(), "redirectAfterLogin")
	if redirectPath == "" {
		redirectPath = "/snippet/create"
	}
	h.app.SessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

func (h *UserHandler) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := h.app.SessionManager.RenewToken(r.Context())
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	h.app.SessionManager.Remove(r.Context(), "authenticatedUserID")
	h.app.SessionManager.Put(r.Context(), "flash", "Logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UserHandler) UserView(w http.ResponseWriter, r *http.Request) {
	td := h.app.NewTemplateData(r)
	if td.User == nil {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	h.app.Render(w, http.StatusOK, "account.tmpl.html", td)
}
