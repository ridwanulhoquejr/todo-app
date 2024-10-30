package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ridwanulhoquejr/todo-app/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// Define a custom ErrDuplicateEmail error.
var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// UserModel for model dependencies in app struct
// so that we can access thes methods (Get, Insert...) from the handlers
type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     password  `json:"-"`
	Activated    bool      `json:"activated"`
	CreationTime time.Time `json:"creation_time"`
}

// Create a custom password type which is a struct containing the plaintext and hashed
// versions of the password for a user. The plaintext field is a *pointer* to a string,
// so that we're able to distinguish between a plaintext password not being present in
// the struct at all, versus a plaintext password which is the empty string "".
type password struct {
	plainPassword *string
	hash          []byte
}

func (p *password) Set(plaintextPassword string) error {

	// hash the plain_password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	// then assign it to the custom password struct
	p.plainPassword = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {

	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))

	// Incorrect pass!
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	// correct pass!
	return true, nil
}

// database methods
func (m *UserModel) Insert(user *User) error {

	query :=
		`
		INSERT INTO 
			users (name, email, password_hash, activated)
			VALUES ($1, $2, $3, $4)
		RETURNING
			id, creation_time
		`
	// args
	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	// create a context for limiting the db-query time limit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// run the db query
	// then Scan the returning values to the user struct
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreationTime,
	)

	// check if there is any error!
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	fmt.Println("in the data -> DB insertion success!")
	return nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {

	query := `
		SELECT 
			id, name, email, password_hash, activated, creation_time
		WHERE 
			email = $1
		LIMIT 1
	`
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.CreationTime,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// validation methods for users
func ValidateEmail(v *validator.Validator, email string) {

	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {

	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 6, "password", "must be at least 6 characters long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {

	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 characters long")

	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)

	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plainPassword != nil {
		ValidatePasswordPlaintext(v, *user.Password.plainPassword)
	}

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
