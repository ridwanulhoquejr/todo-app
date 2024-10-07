package data

import "database/sql"

type Models struct {
	Todo TodoModel
	User UserModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Todo: TodoModel{DB: db},
		User: UserModel{DB: db},
	}
}
