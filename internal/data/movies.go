package data

import (
	"database/sql"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`                       // Unique integer ID for movies
	CreatedAt time.Time `json:"-"`                        // Timestamp for creation of a movie
	Title     string    `json:"title"`                    // String title for movie
	Year      int32     `json:"year,omitempty"`           // Movie release year
	Runtime   Runtime   `json:"runtime,omitempty,string"` // Movie time in minutes
	Genres    []string  `json:"genres,omitempty"`         // Slice of genres for the movie
	Version   int32     `json:"version"`                  // The version number starts at 1 and will be incremented each time the movie information is updated
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (movieModel MovieModel) Insert(movie *Movie) error {
	return nil
}

// Add a placeholder method for fetching a specific record from the movies table.
func (movieModel MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (movieModel MovieModel) Update(movie *Movie) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (movieModel MovieModel) Delete(id int64) error {
	return nil
}
