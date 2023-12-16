package application

import (
	"Meow/internal/data"
	Validator "Meow/internal/validator"
	"fmt"
	"net/http"
	"time"
)

// Add createMovieHandler for the "POST /v1/movies" endpoint.
func (app *Application) createNewMovieHandler(writer http.ResponseWriter, request *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// Using readJSON() helper to decode and also get better error message
	// if any error exist in decoding.
	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	// Initialize a new Validator instance.
	validator := Validator.New()

	// Use the Check() method to execute our validation checks.

	// Title
	validator.Check(input.Title != "", "title", "must be provided")
	validator.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	// Year
	validator.Check(input.Year != 0, "year", "must be provided")
	validator.Check(input.Year >= 1800, "year", "must be greater than 1888")
	validator.Check(input.Year <= int32(time.Now().Year()), "year", "invalid year")

	// Runtime
	validator.Check(input.Runtime != 0, "runtime", "must be provided")
	validator.Check(input.Runtime > 0, "runtime", "must be a positive integer")

	// Genres
	validator.Check(input.Genres != nil, "genres", "must be provided")
	validator.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	validator.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	validator.Check(Validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	// Use the Valid() method to see if any of the checks failed.
	if !validator.Valid() {
		app.failedValidationResponse(writer, request, validator.Errors)
		return
	}

	fmt.Fprintf(writer, "%+v\n", input)
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint.
func (app *Application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {

	// Get id from readIdParam helper.
	id, err := app.readIdParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	// Create a movie with dummy data.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Hello :-)",
		Runtime:   180,
		Genres:    []string{"one", "two", "three"},
		Version:   1,
	}

	// Encoding struct to json and write it in HTTP response body.
	err = app.writeJSON(writer, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
