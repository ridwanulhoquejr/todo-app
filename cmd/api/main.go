package main

import (
	"log"
	"os"
)

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config *config
	logger *log.Logger
}

func main() {

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "Todo API: ", log.Ldate|log.Ltime)

	// initialize a config var
	var cfg config

	app := &application{
		config: cfg.configs(),
		logger: logger,
	}
	err := app.serve()

	if err != nil {
		logger.Print(err)
		return
	}

}
