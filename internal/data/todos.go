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

func (m *TodoModel) Insert(user *User) error {
	return nil
}
func (m *TodoModel) Get(id int64) error {
	return nil
}
func (m *TodoModel) GetAll() error {
	return nil
}
func (m *TodoModel) Update(user *User) error {
	return nil
}
func (m *TodoModel) Delete(id int64) error {
	return nil
}
