package data

import (
	"database/sql"
	"errors"
)

var (
	// Define a custom ErrRecordNotFound error.
	ErrRecordNotFound = errors.New("record not found")

	// Define ErrEditConflict and return it when two threads try change same data at same time.
	ErrEditConflict = errors.New("edit conflict")
)

// Create a model struct which wraps the MovieModel.
type Models struct {
	Movies interface {
		Insert(*Movie) error
		Get(int64) (*Movie, error)
		Update(*Movie) error
		Delete(int64) error
		GetAll(string, []string, Filters) ([]*Movie, error)
	}
}

// Add a New() method which returns a Models struct
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{
			DB: db,
		},
	}
}

// Add NewMockModels()  which returns a Models instance containing the mock models
// only.
func NewMockModels() Models {
	return Models{
		Movies: MockMovieModel{},
	}
}
