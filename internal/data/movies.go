package data

import "time"

type Movie struct {
	ID        int64     // Unique integer ID for movies
	CreatedAt time.Time // Timestamp for creation of a movie
	Title     string    // String title for movie
	Year      int32     // Movie release year
	Runtime   int32     // Movie time in minutes
	Genres    []string  // Slice of genres for the movie
	Version   int32     // The version number starts at 1 and will be incremented each time the movie information is updated
}
