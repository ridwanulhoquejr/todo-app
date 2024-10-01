package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() chi.Router {

	r := chi.NewRouter()

	r.Get("/v1/healthcheck", app.healthcheckHandler)

	// group route
	r.Route("/todo", func(r chi.Router) {
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from the group route\n"))
		})
	})

	return r
}
