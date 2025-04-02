package main

import (
	"errors"
	"net/http"

	"github.com/ridwanulhoquejr/todo-app/internal/data"
	"github.com/ridwanulhoquejr/todo-app/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Convert JSON to GO struct
	// process the payload for creating an User!
	// for that, we need to extract that payload into a go-struct

	var payload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// use our readJSON method for converting the JSON
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.PrintInfo("USER PAYLOAD Parsed", map[string]string{
		"name":  payload.Name,
		"email": payload.Email,
	})

	// 2. Payload to User (entity) convertion
	// Copy the data from the request body into a new User struct.
	user := data.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Activated: false,
	}
	// process the password, convert the plain-password to hash
	err = user.Password.Set(payload.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	app.logger.PrintInfo("USER Password", map[string]string{
		"password": payload.Password,
	})

	// 3. Validation
	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, &user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	app.logger.PrintInfo("USER Validation", map[string]string{
		"validation": "success",
	})

	// 4. perform the db transactions
	err = app.models.User.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 5. return the response!
	err = app.writeJSON(w, http.StatusCreated, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
