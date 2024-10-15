package main

import (
	"os"

	"github.com/ridwanulhoquejr/todo-app/internal/data"
	"github.com/ridwanulhoquejr/todo-app/internal/db"
	"github.com/ridwanulhoquejr/todo-app/internal/jsonlog"
)

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config *config
	logger *jsonlog.Logger
	models *data.Models
}

func main() {

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// db connection
	db, err := db.NewDatabase()
	if err != nil {
		// Use the PrintFatal() method to write a log entry containing the error at the
		// FATAL level and exit. We have no additional properties to include in the log
		// entry, so we pass nil as the second parameter.
		logger.PrintFatal(err, nil)
		return
	}
	defer db.DB.Close()

	// Likewise use the PrintInfo() method to write a message at the INFO level.
	logger.PrintInfo("database connection pool established", nil)

	// if err := db.MigrateDB(); err != nil {
	// 	fmt.Printf("Failed to migrate the database: %w", err)
	// 	return
	// }
	// fmt.Println("Succesfully ping the database")

	app := &application{
		config: Configs(),
		logger: logger,
		models: data.NewModels(db.DB),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
		return
	}
}
