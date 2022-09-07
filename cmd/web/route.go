package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {

	//dynamicMiddleware := alice.New(app.session.Enable, app.authenticatedUser)
	router := mux.NewRouter()
	router.Use(app.session.Enable)

	snippetRouter := router.PathPrefix("/snippet").Subrouter()
	snippetRouter.Use(app.authenticatedUser)

	router.HandleFunc("/", app.home).Methods(http.MethodGet)

	snippetRouter.HandleFunc("/create", app.showSnippetCreationForm).Methods(http.MethodGet)
	snippetRouter.HandleFunc("/create", app.createSnippet).Methods(http.MethodPost)
	snippetRouter.HandleFunc("/{uuid}", app.showSnippet).Methods(http.MethodGet)
	snippetRouter.HandleFunc("/delete/{uuid}", app.deleteSnippet).Methods(http.MethodGet)

	router.HandleFunc("/user/signup", app.showSignupForm).Methods(http.MethodGet)
	router.HandleFunc("/user/signup", app.signupUser).Methods(http.MethodPost)
	router.HandleFunc("/user/login", app.showLoginForm).Methods(http.MethodGet)
	router.HandleFunc("/user/login", app.loginUser).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))
	return router
}
