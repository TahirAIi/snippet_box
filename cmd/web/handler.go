package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"snippet_box/pkg/forms"
	"snippet_box/pkg/models"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
	}

	app.render(writer, request, "home.page.tmpl", &templateData{Snippets: snippets})
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
		app.render(writer, request, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))

	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	snippetId, _ := strconv.Atoi(params["id"])

	if snippetId < 1 {
		app.notFound(writer)
		return
	}
	snippet, err := app.snippets.Get(snippetId)
	if err == models.ErrorNoRecord {
		app.notFound(writer)
		return
	} else if err != nil {
		app.serverError(writer, err)
		return
	}
	app.render(writer, request, "show.page.tmpl", &templateData{Snippet: snippet})
}

func (app *application) deleteSnippet(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	snippetId, _ := strconv.Atoi(params["id"])
	if snippetId < 1 {
		app.notFound(writer)
		return
	}
	app.snippets.Delete(snippetId)
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
