package application

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Retrive the "id" parameter from the current request context.
// convert it from string to int and return it.
func (app *Application) readIdParam(request *http.Request) (int64, error) {

	// When httprouter parsing a request, any URL parameters will be stored
	// in the request context. We can use the ParamsFromContext() function to
	// retrive a slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(request.Context())

	// We can use ByName() method to get the value of the "id" parameter from the slice.
	// In our project all movies have an id promary key. But when we get "id" parameter from request context
	// it will be string instead of int. We need to convert and validate it.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	// Return id for id and nil for error.
	return id, nil
}
