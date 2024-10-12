package main

import (
	"log"
	"os"

	"github.com/ridwanulhoquejr/todo-app/internal/data"
	"github.com/ridwanulhoquejr/todo-app/internal/db"
)

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config *config
	logger *log.Logger
	models *data.Models
}

func main() {

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "Todo API: ", log.Ldate|log.Ltime)

	// db connection
	db, err := db.NewDatabase()
	if err != nil {
		logger.Fatal("Failed to connect with database: %w", err)
		return
	}
	defer db.DB.Close()

	// if err := db.MigrateDB(); err != nil {
	// 	fmt.Printf("Failed to migrate the database: %w", err)
	// 	return
	// }
	// fmt.Println("Succesfully ping the database")

	// initialize a config var
	// var cfg config

	app := &application{
		config: Configs(),
		logger: logger,
		models: data.NewModels(db.DB),
	}

	err = app.serve()

	if err != nil {
		logger.Print(err)
		return
	}

}
