package application

import (
	"Meow/internal/data"
	"fmt"
	"net/http"
	"time"
)

// Add createMovieHandler for the "POST /v1/movies" endpoint.
func (app *Application) createNewMovieHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Creating new movie.")
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint.
func (app *Application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {

	// Get id from readIdParam helper.
	id, err := app.readIdParam(request)
	if err != nil {
		http.NotFound(writer, request)
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
	err = app.writeJSON(writer, http.StatusOK, movie, nil)
	if err != nil {
		app.Logger.Println(err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}
}
