package main

import (
	"net/http"

	"github.com/ridwanulhoquejr/todo-app/internal/data"
	"github.com/ridwanulhoquejr/todo-app/internal/validator"
)

func (app *application) singleTodoHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"id": id}, nil)

}

func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. extract the payload
	var payload struct {
		Title      string `json:"title"`
		Descripton string `json:"description"`
		UserID     int64  `json:"user_id"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// 2. payload to Entity
	todo := data.Todo{
		Title:      payload.Title,
		Descripton: payload.Descripton,
		UserID:     payload.UserID, // this would be coming from token
	}

	// 3. validaton
	v := validator.New()

	//! should work with user_id validaiton
	if data.ValidateTodo(v, &todo); !v.Valid() {
		app.logger.Printf("validation error: %s", v.Errors)
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 4. communicate with db
	err = app.models.Todo.Insert(&todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// 5. return response
	err = app.writeJSON(w, http.StatusCreated, todo, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	return
}
