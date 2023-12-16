package validator

import (
	"Meow/internal/data"
	"time"
)

func MovieValidator(v *Validator, movie *data.Movie) {
	// Use the Check() method to execute our validation checks.
	// Title
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	// Year
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1800, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "invalid year")

	// Runtime
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	// Genres
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(Unique(movie.Genres), "genres", "must not contain duplicate values")
}
