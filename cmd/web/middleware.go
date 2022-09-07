package main

import (
	"net/http"
)

func (app *application) authenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !app.isUserAuthenticated(request) {
			http.Redirect(writer, request, "/user/login", http.StatusFound)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
