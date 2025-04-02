package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	// dsn := fmt.Sprintf(
	// 	"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// 	os.Getenv("POSTGRES_USER"),
	// 	os.Getenv("POSTGRES_DB"),
	// 	os.Getenv("POSTGRES_PASSWORD"),
	// 	os.Getenv("SSL_MODE"),
	// )
	// fmt.Printf("db connection string: %s", dsn)

	connStr := "user=todo dbname=todo_db password=root sslmode=disable"

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return &Database{}, fmt.Errorf("could not cnnect to the databse: %w", err)
	}

	fmt.Println("successfully connected to the db")
	return &Database{
		DB: dbConn,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {

	return d.DB.PingContext(ctx)
}
