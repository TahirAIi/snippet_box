package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/", app.home).Methods(http.MethodGet)
	router.HandleFunc("/snippet/create", app.showSnippetCreationForm).Methods(http.MethodGet)
	router.HandleFunc("/snippet/create", app.createSnippet).Methods(http.MethodPost)
	router.HandleFunc("/snippet/{id:[0-9]+}", app.showSnippet).Methods(http.MethodGet)
	router.HandleFunc("/snippet/delete/{id:[0-9]+}", app.deleteSnippet).Methods(http.MethodGet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))
	return router
}
