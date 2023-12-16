package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Define an envelope type.
type envelope map[string]interface{}

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
func (app *Application) writeJSON(writer http.ResponseWriter, status int, value envelope, headers http.Header) error {
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

// A function for reading json content which send from user and detect
// its errors and replace our error message instead of defualt http errors.
func (app *Application) readJSON(writer http.ResponseWriter, request *http.Request, dest interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	// Decode the request body into destination.
	err := decoder.Decode(&dest)
	if err != nil {
		// If there is an error during decoding...
		var syntaxErr *json.SyntaxError
		var unMarshalTypeErr *json.UnmarshalTypeError
		var invalidUnMarshalErr *json.InvalidUnmarshalError
		switch {

		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError.
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxErr.Offset)

		// In some cases Decode() may also return an io.ErrUnexpectedEOF
		// error for syntax error in JSON.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unMarshalTypeErr):
			if unMarshalTypeErr.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unMarshalTypeErr.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unMarshalTypeErr.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body can not be empty")

		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"".
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", field)

		// If the request body exceeds 1MB in size the decode will now fail with the
		// error "http: request body too large".
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic.
		case errors.As(err, &invalidUnMarshalErr):
			panic(err)

		// For other cases, return defualt error message.
		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination.
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	// Decoding finished without any error.
	return nil
}
