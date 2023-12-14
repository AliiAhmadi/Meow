package application

import (
	"encoding/json"
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

// Define writeJSON() function for json encoding and write in response body.
func (app *Application) writeJSON(writer http.ResponseWriter, status int, value interface{}, headers http.Header) error {
	// Encoding data to json and check for errors.
	js, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	for key, value := range headers {
		writer.Header()[key] = value
	}

	// Add "Content-Type: application/json" header. Then write
	// status code.
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(js)

	return nil
}
