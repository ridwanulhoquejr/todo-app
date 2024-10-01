package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	app.logger.Printf("starting %s server on port %s", app.config.env, srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err)
		return err
	}

	return nil
}
