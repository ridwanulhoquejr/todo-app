package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ridwanulhoquejr/todo-app/internal/data"
	"github.com/ridwanulhoquejr/todo-app/internal/validator"
)

func (app *application) delelteTokenHandler(w http.ResponseWriter, r *http.Request) {

}
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

	// extract the payload
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// process it with readJSON helper method
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate payload
	v := validator.New()

	data.ValidateEmail(v, payload.Email)
	data.ValidatePasswordPlaintext(v, payload.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// check with user table, if email is exist or not
	user, err := app.models.User.GetByEmail(payload.Email)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// at this moment we sure that user is found!
	// so, check the password.Hash
	// if it is matches; then generate a token
	match, err := user.Password.Matches(payload.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// generate a token
	token, err := app.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
