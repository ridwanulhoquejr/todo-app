package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/ridwanulhoquejr/todo-app/internal/validator"
)

type Todo struct {
	Title        string    `json:"title"`
	Descripton   string    `json:"description"`
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Completed    bool      `json:"completed"`
	Version      int32     `json:"verison"`
	CreationTime time.Time `json:"creation_time"`
}

type TodoModel struct {
	DB *sql.DB
}

// here Todo struct methods will communicate with the Database

func (m *TodoModel) Insert(td *Todo) error {
	// 1. query
	query :=
		`
		INSERT INTO 
			todo (title, description, user_id)
			values ($1, $2, $3)
		RETURNING
			id, completed, version, creation_time
		`
	// 2. args
	args := []any{td.Title, td.Descripton, td.UserID}

	// 3. create a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&td.ID,
		&td.Completed,
		&td.Version,
		&td.CreationTime,
	)
	if err != nil {
		return err
	}

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

func ValidateTodo(v *validator.Validator, todo *Todo) {
	v.Check(todo.Title != "", "title", "must be provided")
	// v.Check(todo.UserID == 0, "user_id", "must be provided")
}
