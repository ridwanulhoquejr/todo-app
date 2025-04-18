package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ridwanulhoquejr/todo-app/internal/data"
	"github.com/ridwanulhoquejr/todo-app/internal/validator"
)

var errInvalidPathParam = errors.New("invalid id: path parameter must be greater than zero")

func (app *application) getQueries(
	r *http.Request, v *validator.Validator) data.Queries {

	qs := r.URL.Query()
	q := data.Queries{}
	// get the limit and Offset and assign it to our Queries.Pagination
	q.Pagination.Page = app.readInt(qs, "page", 1, v)
	q.Pagination.PageSize = app.readInt(qs, "page_size", 5, v)

	// search query
	q.Search.Title = app.readString(qs, "title", "")
	// get the Filters
	q.Sorts.Sort = app.readString(qs, "sort", "-creation_time")
	// add a Sortsafelist
	q.Sorts.SafeList = []string{"title", "id", "creation_time", "-title", "-id", "-creation_time"}

	// get time for range filters
	now := time.Now()
	oneMonthAgo := now.AddDate(0, -1, 0)
	q.Filters.StartDate = app.readTime(qs, "start_date", oneMonthAgo)
	q.Filters.EndDate = app.readTime(qs, "end_date", now)

	return q
}

func (app *application) getTodoHandler(w http.ResponseWriter, r *http.Request) {

	// extract the user info from r.Context()
	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if id <= 0 {
		app.badRequestResponse(w, r, errInvalidPathParam)
		return
	}

	todo, err := app.models.Todo.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, todo, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllTodoHandler(w http.ResponseWriter, r *http.Request) {

	// extract the user info from r.Context()
	user := app.contextGetUser(r)

	fmt.Printf("--> User from authenticate middleware context: %+v", user)

	v := validator.New()
	queries := app.getQueries(r, v)

	// execute validate check
	if data.ValidateQueries(v, queries); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	todos, metadata, err := app.models.Todo.GetAll(user.ID, queries)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		app.logger.PrintError(err, map[string]string{"todos": "error returned from GetAll db"})
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todos, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
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

	if data.ValidateTodo(v, &todo); !v.Valid() {
		app.logger.PrintError(err, map[string]string{"validation": "error returned from ValidateTodo"})
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 4. communicate with db
	err = app.models.Todo.Insert(&todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/todos/%d", todo.ID))

	// 5. return response
	err = app.writeJSON(w, http.StatusCreated, todo, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {

	// get the user from the auth_context
	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	if id <= 0 {
		app.badRequestResponse(w, r, errInvalidPathParam)
		return
	}

	// call the db method
	// TODO: user_id should come from the AUTHENTICATIONS!
	err = app.models.Todo.Delete(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// return the response!
	err = app.writeJSON(w, http.StatusOK, "todo succesfully deleted", nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {

	// get the user from the auth_context
	user := app.contextGetUser(r)

	// get the id of the Todo
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if id <= 0 {
		app.badRequestResponse(w, r, errInvalidPathParam)
		return
	}

	// retrieve the Todo from DB using that id
	// TODO: user_id should come from the AUTHENTICATIONS!
	todo, err := app.models.Todo.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// define a payload struct
	var payload struct {
		Title      *string `json:"title"`
		Descripton *string `json:"description"`
		Completed  *bool   `json:"completed"`
	}

	// extract the request payload to our defined payload
	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// as our payload is a pointers, we can check if it is nil
	if payload.Title != nil {
		todo.Title = *payload.Title
	}
	if payload.Descripton != nil {
		todo.Descripton = *payload.Descripton
	}
	if payload.Completed != nil {
		todo.Completed = *payload.Completed
	}

	// validation
	v := validator.New()
	v.Check(todo.Title != "", "title", "must be provided")

	if !v.Valid() {
		app.logger.PrintError(err, map[string]string{"todo-update": "error returned from update-todo validator"})
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// call the db method
	// TODO: user_id should come from the AUTHENTICATIONS!
	err = app.models.Todo.Update(user.ID, todo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// return back the responses
	err = app.writeJSON(w, http.StatusOK, todo, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
