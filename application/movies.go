package application

import (
	"fmt"
	"net/http"
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

	fmt.Fprintf(writer, "here is movie with %d id", id)
}
