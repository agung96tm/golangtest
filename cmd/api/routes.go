package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/login", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/users", app.listUserHandler)
	router.HandlerFunc(http.MethodPost, "/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPatch, "/users/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/users/:id", app.deleteUserHandler)

	return router
}
