package db

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const maxRetries = 5
const retryDelay = 2 * time.Second

func (d *Database) MigrateDB() error {
	// Check if the DB connection is nil
	if d.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	fmt.Println("Migrating the database")

	// Retry logic for connecting to the database
	var driver database.Driver
	var err error

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Trying for %d times:", i)
		driver, err = postgres.WithInstance(d.DB, &postgres.Config{})
		if err == nil {
			break
		}
		fmt.Printf("Failed to create driver, retrying in %v...\n", retryDelay)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return fmt.Errorf("could not create the postgres driver after %d retries: %w", maxRetries, err)
	}

	// Initialize migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not initialize migration instance: %w", err)
	}

	// Apply all the migrations (up)
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("No new migrations to apply")
	} else {
		fmt.Println("Successfully migrated the database")
	}

	return nil
}
