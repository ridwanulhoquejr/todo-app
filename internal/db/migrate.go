package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func (d *Database) MigrateDB() error {
	// Check if the DB connection is nil
	if d.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	fmt.Println("Migrating the database")

	// Create the PostgreSQL driver for migration
	driver, err := postgres.WithInstance(d.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %w", err)
	}

	// Point to the correct file path for your migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations", // Adjust this path based on your folder structure
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

func (d *Database) RollbackDB() error {
	// Check if the DB connection is nil
	if d.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	fmt.Println("Rolling back the database")

	// Create the PostgreSQL driver for migration
	driver, err := postgres.WithInstance(d.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %w", err)
	}

	// Initialize the migrate instance with the correct file path
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations", // Adjust this path as per your file structure
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not initialize migration instance: %w", err)
	}

	// Rollback all migrations (full down)
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not roll back migrations: %w", err)
	}

	fmt.Println("Successfully rolled back the database")
	return nil
}
