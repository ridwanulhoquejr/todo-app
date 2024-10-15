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

	// Again, we use the PrintInfo() method to write a "starting server" message at the
	// INFO level. But this time we pass a map containing additional properties (the
	// operating environment and server address) as the final parameter.
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  Configs().env,
	})

	err := srv.ListenAndServe()
	if err != nil {
		app.logger.PrintFatal(err, nil)
		return err
	}
	return nil
}
