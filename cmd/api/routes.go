package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() chi.Router {

	r := chi.NewRouter()
	// override built-in NotFound handler
	r.NotFound(app.URLNotFound)

	r.Get("/v1/healthcheck", app.healthcheckHandler)

	// todo routes group
	r.Route("/todos", func(r chi.Router) {
		r.Get("/{id}", app.singleTodoHandler)
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from GET ALL of todos\n"))
		})
		// r.Post("", app.createTodoHandler)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, from POST Method of Todos\n"))
		})
		// r.Patch("", app.updateTodoHandler)
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, from PATCH Method of Todos\n"))

		})
		// r.Delete("/{id}", app.deleteTodoHandler)
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, from DELETE Method of Todos\n"))

		})

	})

	// user route group
	r.Group(func(r chi.Router) {
		r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, from GET Method of Users\n"))
		})
	})

	return r
}
