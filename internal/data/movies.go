package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`         // Unique integer ID for movies
	CreatedAt time.Time `json:"created_at"` // Timestamp for creation of a movie
	Title     string    `json:"title"`      // String title for movie
	Year      int32     `json:"year"`       // Movie release year
	Runtime   int32     `json:"runtime"`    // Movie time in minutes
	Genres    []string  `json:"genres"`     // Slice of genres for the movie
	Version   int32     `json:"version"`    // The version number starts at 1 and will be incremented each time the movie information is updated
}
