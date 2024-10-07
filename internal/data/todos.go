package data

import (
	"database/sql"
	"time"
)

type Todo struct {
	Title        string    `json:"title"`
	Descripton   string    `json:"description"`
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Completed    bool      `json:"completed"`
	Version      int32     `json:"verison"`
	CreationTime time.Time `json:"creation_time,-"`
}

type TodoModel struct {
	DB *sql.DB
}

// here Todo struct methods will communicate with the Database
