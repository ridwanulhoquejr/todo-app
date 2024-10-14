package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
func (m *TodoModel) Insert(todo *Todo) error {
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
	args := []any{todo.Title, todo.Descripton, todo.UserID}

	// 3. create a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&todo.ID,
		&todo.Completed,
		&todo.Version,
		&todo.CreationTime,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *TodoModel) Get(id, userId int64) (*Todo, error) {

	// 1. write the query
	query :=
		`
	SELECT id, title, description, completed, creation_time, version 
	FROM todo
		WHERE id = $1 
			AND user_id = $2
	ORDER BY creation_time DESC
	LIMIT 1
	`
	args := []any{id, userId}

	var todo Todo

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Descripton,
		&todo.Completed,
		&todo.CreationTime,
		&todo.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &todo, nil
}

func (m *TodoModel) GetAll(userId int64, q Queries) ([]*Todo, error) {
	query :=
		fmt.Sprintf(`
				SELECT
					id, user_id, title, description, completed, creation_time, version
				FROM todo
					WHERE
						user_id = $1
						AND (to_tsvector('simple', title) @@ plainto_tsquery('simple', $2) OR $2 = '')
						AND (creation_time BETWEEN $3 AND $4)
				ORDER BY %s %s, id ASC
				LIMIT $5 OFFSET $6`,
			q.Sorts.sortColumn(), q.Sorts.sortDirection())

	args := []any{userId, q.Search.Title, q.Filters.StartDate, q.Filters.EndDate, q.Pagination.limit(), q.Pagination.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []*Todo{}

	for rows.Next() {
		var todo Todo

		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Descripton,
			&todo.Completed,
			&todo.CreationTime,
			&todo.Version,
		)
		// handle the Scan err
		if err != nil {
			return nil, err
		}

		todos = append(todos, &todo)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err := rows.Err(); err != nil {
		return nil, nil
	}

	return todos, nil
}

func (m *TodoModel) Update(userID int64, todo *Todo) error {
	query := `
		UPDATE todo
			SET
				title = $1,
				description = $2,
				completed = $3,
				version = version +1
			WHERE id = $4
				AND user_id = $5
				AND version = $6
		RETURNING version
	`
	args := []any{todo.Title, todo.Descripton, todo.Completed, todo.ID, userID, todo.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m *TodoModel) Delete(id, userID int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE from todo WHERE id = $1 and user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateTodo(v *validator.Validator, todo *Todo) {
	v.Check(todo.Title != "", "title", "must be provided")
	v.Check(todo.UserID != 0, "user_id", "must be provided")
}
