package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"snippet_box/pkg/forms"
	"snippet_box/pkg/models"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
	}

	app.render(writer, request, "home.page.tmpl", &templateData{Snippets: snippets, AuthenticatedUser: app.isUserAuthenticated(request)})
}

func (app *application) showSnippetCreationForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "create.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusUnprocessableEntity)
		return
	}

	form := forms.New(request.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(writer, request, "create.page.tmpl", &templateData{
			Form:              form,
			AuthenticatedUser: app.isUserAuthenticated(request),
		})
		return
	}

	uuid, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))

	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/%s", uuid), http.StatusSeeOther)
}

func (app *application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	snippetUuid, _ := params["uuid"]
	_, err := uuid.Parse(snippetUuid)
	if err != nil {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}
	snippet, err := app.snippets.Get(snippetUuid)
	if err == models.ErrorNoRecord {
		fmt.Println("No record found in db")
		return
	} else if err != nil {
		app.serverError(writer, err)
		return
	}
	app.render(writer, request, "show.page.tmpl", &templateData{
		Snippet:           snippet,
		AuthenticatedUser: app.isUserAuthenticated(request),
	})
}

func (app *application) deleteSnippet(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	snippetId, _ := params["uuid"]
	app.snippets.Delete(snippetId)
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func (app *application) showSignupForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "signup.page.tmpl", &templateData{
		Form:              forms.New(nil),
		AuthenticatedUser: app.isUserAuthenticated(request),
	})
}

func (app *application) signupUser(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusUnprocessableEntity)
		return
	}

	form := forms.New(request.Form)
	form.Required("name", "email", "password")
	form.MaxLength("name", 64)
	form.MaxLength("email", 255)

	if !form.Valid() {
		app.render(writer, request, "signup.page.tmpl", &templateData{Form: form})
	}
	err = app.users.IsEmailAlreadyTaken(form.Get("email"))
	if err == gorm.ErrRecordNotFound {
		err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
		if err != nil {
			app.serverError(writer, err)
		}
	} else if err == gorm.ErrRegistered {
		form.Errors.Add("email", "This email is already taken")
		app.render(writer, request, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func (app *application) showLoginForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "login.page.tmpl", &templateData{
		Form:              forms.New(nil),
		AuthenticatedUser: app.isUserAuthenticated(request),
	})
}

func (app *application) loginUser(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusUnprocessableEntity)
		return
	}
	form := forms.New(request.Form)
	form.Required("email", "password")

	if !form.Valid() {
		app.render(writer, request, "login.page.tmpl", &templateData{
			Form:              form,
			AuthenticatedUser: app.isUserAuthenticated(request),
		})
	}
	user, err := app.users.Login(form.Get("email"), form.Get("password"))
	if err != nil {
		app.serverError(writer, err)
	}
	app.session.Put(request, "userId", user.Uuid)
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
