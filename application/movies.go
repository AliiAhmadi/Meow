package application

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Add createMovieHandler for the "POST /v1/movies" endpoint.
func (app *Application) createNewMovieHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Creating new movie.")
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint.
func (app *Application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {

	// When httprouter parsing a request, any URL parameters will be stored
	// in the request context. We can use the ParamsFromContext() function to
	// retrive a slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(request.Context())

	// We can use ByName() method to get the value of the "id" parameter from the slice.
	// In our project all movies have an id promary key. But when we get "id" parameter from request context
	// it will be string instead of int. We need to convert and validate it.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "here is movie with %d id", id)
}
