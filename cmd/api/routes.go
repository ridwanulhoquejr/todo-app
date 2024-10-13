package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.Timeout(30 * time.Second))
	// override built-in NotFound handler
	r.NotFound(app.URLNotFound)
	r.Get("/v1/healthcheck", app.healthcheckHandler)

	// todo routes group
	r.Route("/todos", func(r chi.Router) {
		r.Get("/all", app.getAllTodoHandler)
		r.Get("/{id}", app.getTodoHandler)
		r.Patch("/{id}", app.updateTodoHandler)
		r.Delete("/{id}", app.deleteTodoHandler)
		r.Post("/", app.createTodoHandler)
	})

	// user route group
	r.Group(func(r chi.Router) {
		r.Post("/users", app.createUserHandler)
		// r.Get("/", app.createUserHandler)
	})

	return r
}
