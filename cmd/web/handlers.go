package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	td := app.newTemplateData(r)
	td.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl.html", td)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	td := app.newTemplateData(r)
	td.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", td)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)
	td.Form = &snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", td)
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field must be between 1 and 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.ValueIn(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		td := app.newTemplateData(r)
		td.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", td)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// ################################
type usersSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)
	td.Form = &usersSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", td)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form usersSignupForm
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRGX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		td := app.newTemplateData(r)
		td.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", td)
		return
	}

	err := app.users.Insert(form.Name, form.Email, form.Password)
	if errors.Is(err, models.ErrDuplicateEmail) {
		form.AddFieldError("email", "Email address already in use")
		td := app.newTemplateData(r)
		td.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", td)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User created successfully")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)
	td.Form = &userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", td)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRGX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		td := app.newTemplateData(r)
		td.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", td)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if errors.Is(err, models.ErrInvalidCredentials) {
		form.AddNonFieldError("Email or password is incorrect")
		td := app.newTemplateData(r)
		td.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", td)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "Logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
