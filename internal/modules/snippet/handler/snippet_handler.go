package handler

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"snippetbox/internal/modules/snippet/domain"
	"snippetbox/internal/modules/snippet/service"
	"snippetbox/internal/platform/database"
	"snippetbox/internal/platform/web"
	"snippetbox/internal/validator"
	"strconv"
)

type SnippetHandler struct {
	service *service.Snippet
	app     *web.Application
}

func NewSnippet(app *web.Application, svc *service.Snippet) *SnippetHandler {
	return &SnippetHandler{
		service: svc,
		app:     app,
	}
}

func (h *SnippetHandler) SnipperList(w http.ResponseWriter, r *http.Request) {
	td := h.app.NewTemplateData(r)
	snippets, err := h.service.SnippetList()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	td.List = snippets
	h.app.Render(w, http.StatusOK, "home.tmpl.html", td)
}

func (h *SnippetHandler) SnippetView(w http.ResponseWriter, r *http.Request) {
	// TODO: remove httprouter dependency
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		h.app.NotFound(w)
		return
	}

	snippet, err := h.service.SnippetRead(id)
	if errors.Is(err, database.ErrNoRecord) {
		h.app.NotFound(w)
		return
	} else if err != nil {
		h.app.ServerError(w, err)
		return
	}

	td := h.app.NewTemplateData(r)
	td.Model = snippet

	h.app.Render(w, http.StatusOK, "view.tmpl.html", td)
}

func (h *SnippetHandler) SnippetCreateView(w http.ResponseWriter, r *http.Request) {
	td := h.app.NewTemplateData(r)
	td.Form = &domain.SnippetCreateForm{Expires: 365}

	h.app.Render(w, http.StatusOK, "create.tmpl.html", td)
}

func (h *SnippetHandler) SnippetCreatePost(w http.ResponseWriter, r *http.Request) {
	userId := h.app.SessionManager.GetInt(r.Context(), "authenticatedUserID")
	var form domain.SnippetCreateForm
	if err := h.app.DecodePostForm(r, &form); err != nil {
		h.app.ClientError(w, http.StatusBadRequest)
		return
	}

	id, err := h.service.SnippetCreate(&form, userId)
	if errors.Is(err, validator.ErrInvalidForm) {
		td := h.app.NewTemplateData(r)
		td.Form = form

		h.app.Render(w, http.StatusUnprocessableEntity, "create.tmpl.html", td)
		return
	} else if err != nil {
		h.app.ServerError(w, err)
		return
	}

	h.app.SessionManager.Put(r.Context(), "flash", "Snippet created successfully")
	http.Redirect(w, r, "/snippet/view/"+strconv.Itoa(id), http.StatusSeeOther)
}
