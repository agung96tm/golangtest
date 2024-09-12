package main

import (
	"errors"
	"github.com/pascaldekloe/jwt"
	"golangtest/internal/data"
	"golangtest/internal/validator"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) listUserHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   users,
	}); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleId   string `json:"role_id"`
		RoleName string `json:"role_name"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.badRequestResponse(w, err)
		return
	}

	v := validator.New()
	v.Check(input.RoleId != "", "role_id", "required")
	v.Check(input.Name != "", "name", "required")
	v.Check(input.Email != "", "email", "required")
	v.Check(input.Password != "", "password", "required")
	if !v.IsValid() {
		app.badValidatorResponse(w, v.Errors)
		return
	}

	err := app.models.Users.Create(&data.User{
		Name:       input.Name,
		Email:      input.Email,
		Password:   input.Password,
		RoleID:     input.RoleId,
		RoleName:   input.RoleName,
		LastAccess: time.Now(),
	})
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "successfully",
	}); err != nil {
	}
}

func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.badRequestResponse(w, err)
		return
	}

	user, err := app.models.Users.GetByID(id)
	if err != nil {
		app.notFoundError(w)
		return
	}

	v := validator.New()
	if input.Name != "" {
		v.Check(input.Name != "", "name", "required")
		user.Name = input.Name
	}
	if !v.IsValid() {
		app.badValidatorResponse(w, v.Errors)
		return
	}

	err = app.models.Users.Update(user)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "successfully",
	}); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	user, err := app.models.Users.GetByID(id)
	if err != nil {
		app.notFoundError(w)
		return
	}
	err = app.models.Users.DeleteByID(user.ID)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if err := app.writeJSON(w, http.StatusNoContent, map[string]interface{}{
		"status":  "true",
		"message": "successfully",
	}); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.badRequestResponse(w, err)
		return
	}

	v := validator.New()
	v.Check(input.Email != "", "email", "required")
	v.Check(input.Password != "", "password", "required")
	if !v.IsValid() {
		app.badValidatorResponse(w, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		app.notFoundError(w)
		return
	}

	isMatched := user.Matches(input.Password)
	if !isMatched {
		app.badRequestResponse(w, errors.New("invalid password"))
		return
	}

	var claims jwt.Claims
	claims.Subject = strconv.Itoa(user.ID)
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.SecretKey))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "successfully",
		"data": map[string]interface{}{
			"access_token": string(jwtBytes),
		},
	}); err != nil {
		app.serverError(w, err)
	}
}
