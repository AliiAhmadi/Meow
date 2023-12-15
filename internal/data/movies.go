package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`                       // Unique integer ID for movies
	CreatedAt time.Time `json:"-"`                        // Timestamp for creation of a movie
	Title     string    `json:"title"`                    // String title for movie
	Year      int32     `json:"year,omitempty"`           // Movie release year
	Runtime   Runtime   `json:"runtime,omitempty,string"` // Movie time in minutes
	Genres    []string  `json:"genres,omitempty"`         // Slice of genres for the movie
	Version   int32     `json:"version"`                  // The version number starts at 1 and will be incremented each time the movie information is updated
}
