package main

import (
	"fmt"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, status int, message interface{}) error {
	var errMessage struct {
		Message interface{} `json:"message"`
	}
	errMessage.Message = message

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	err := app.writeJSON(w, status, errMessage)

	if err != nil {
		_ = app.writeJSON(w, http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	app.writeJSON(w, http.StatusInternalServerError, err.Error())
}

func (app *Application) notFoundError(w http.ResponseWriter) {
	message := "not found"
	app.writeJSON(w, http.StatusNotFound, message)
}

func (app *Application) badRequestResponse(w http.ResponseWriter, err error) {
	message := fmt.Sprintf("invalid request %s", err.Error())
	app.writeJSON(w, http.StatusBadRequest, message)
}

func (app *Application) badValidatorResponse(w http.ResponseWriter, err map[string]string) {
	app.writeJSON(w, http.StatusBadRequest, err)
}
