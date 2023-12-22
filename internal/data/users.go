package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Define custom error
var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// Define a User struct to represent an individual user.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

// Define Queries
var (
	INSERT_USER_QUERY = `INSERT INTO users (name, email, password_hash, activated) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	GET_USER_BY_EMAIL = `SELECT id, created_at, name, email, password_hash, activated, version
	FROM users WHERE email = $1
	`

	UPDATE_USER_QUERY = `UPDATE users SET
	name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version
	`
)

// Create a UserModel struct which wraps the connection pool.
type UserModel struct {
	DB *sql.DB
}

// Insert a new record in the database for the user. Note that the id, created_at and
// version fields are all automatically generated by our database, so we use the
// RETURNING clause to read them into the User struct after the insert, in the same way
// that we did when creating a movie.
func (userModel UserModel) Insert(user *User) error {
	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter.
	err := userModel.DB.QueryRowContext(ctx, INSERT_USER_QUERY, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

// Retrieve the User details from the database based on the user's email address.
// Because we have a UNIQUE constraint on the email column, this SQL query will only
// return one record (or none at all, in which case we return a ErrRecordNotFound error).
func (userModel UserModel) GetByEmail(email string) (*User, error) {
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := userModel.DB.QueryRowContext(ctx, GET_USER_BY_EMAIL, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
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

func (userModel UserModel) Update(user *User) error {
	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := userModel.DB.QueryRowContext(ctx, UPDATE_USER_QUERY, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

// Create a custom password type which is a struct containing the plaintext and hashed
// versions of the password for a user.
type password struct {
	Plaintext *string
	Hash      []byte
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
