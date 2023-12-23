package application

import (
	"fmt"
	"net/http"
)

// The logError() method is a generic helper for logging an error message.
func (app *Application) logError(request *http.Request, err error) {
	app.Logger.PrintError(err, map[string]string{
		"request_method": request.Method,
		"request_url":    request.URL.String(),
	})
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using an interface{}
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func (app *Application) errorResponse(writer http.ResponseWriter, request *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := app.writeJSON(writer, status, env, nil)
	if err != nil {
		app.logError(request, err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

// The serverErrorResponse() method will be used when our application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses the
// errorResponse() helper to send a 500 Internal Server Error status code and JSON
// response (containing a generic error message) to the client.
func (app *Application) serverErrorResponse(writer http.ResponseWriter, request *http.Request, err error) {
	app.logError(request, err)

	app.errorResponse(writer, request, http.StatusInternalServerError, "Internal server error")
}

// The notFoundResponse() method will be used to send a 404 Not Found status code and
// JSON response to the client.
func (app *Application) notFoundResponse(writer http.ResponseWriter, request *http.Request) {
	app.errorResponse(writer, request, http.StatusNotFound, "the requested resource could not be found")
}

// The methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func (app *Application) methodNotAllowedResponse(writer http.ResponseWriter, request *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", request.Method)
	app.errorResponse(writer, request, http.StatusMethodNotAllowed, message)
}

// The badRequestResponse() method will be used when need to send
// error response when invalid request recived.
func (app *Application) badRequestResponse(writer http.ResponseWriter, request *http.Request, err error) {
	app.errorResponse(writer, request, http.StatusBadRequest, err.Error())
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in our Validator type.
func (app *Application) failedValidationResponse(writer http.ResponseWriter, request *http.Request, errors map[string]string) {
	app.errorResponse(writer, request, http.StatusUnprocessableEntity, errors)
}

// When a conflict occures in write data in database
// we send a conflict error with 409 status code.
func (app *Application) editConflictResponse(writer http.ResponseWriter, request *http.Request) {
	app.errorResponse(writer, request, http.StatusConflict, "unable to update the record, please try again")
}

// When rate limit exceeded we will send a
// response error with 429 to many request for user.
func (app *Application) rateLimitExceededResponse(writer http.ResponseWriter, request *http.Request) {
	app.errorResponse(writer, request, http.StatusTooManyRequests, "rate limit exceeded")
}

// Invalid credentials error when user credentials incorrect or user email not found
func (app *Application) invalidCredentialsResponse(writer http.ResponseWriter, request *http.Request) {
	app.errorResponse(writer, request, http.StatusUnauthorized, "invalid authntication credentials")
}

// invalidAuthenticationTokenResponse
func (app *Application) invalidAuthenticationTokenResponse(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("WWW-Authenticate", "Bearer")

	app.errorResponse(writer, request, http.StatusUnauthorized, "invalid or missing authentication token")
}
