package handler

import (
	"errors"
	"net/http"
	"snippetbox/internal/platform/web"
	"snippetbox/internal/user/domain"
	"snippetbox/internal/user/repository"
	"snippetbox/internal/user/service"
)

type UserHandler struct {
	service *service.UserService
	app     *web.Application
}

func NewUserHandler(app *web.Application, svc *service.UserService) *UserHandler {
	return &UserHandler{
		service: svc,
		app:     app,
	}
}

func (h *UserHandler) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	var form domain.UsersSignupForm
	if err := h.app.DecodePostForm(r, &form); err != nil {
		h.app.ClientError(w, http.StatusBadRequest)
		return
	}
	ok, err := h.service.UserSignup(&form)
	if !ok {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use")
		}

		td := h.app.NewTemplateData(r)
		td.Form = form
		h.app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", td)
		return
	}
}
