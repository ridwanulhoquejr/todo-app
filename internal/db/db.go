package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return &Database{}, fmt.Errorf("could not cnnect to the databse: %w", err)
	}
	return &Database{
		DB: dbConn,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}