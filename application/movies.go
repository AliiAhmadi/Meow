package application

import (
	"Meow/internal/data"
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
