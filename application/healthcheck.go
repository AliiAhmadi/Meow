package application

import (
	"fmt"
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *Application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {

	// Create a json response for health check request.
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.Config.Env, app.Version)

	// Set the "Content-Type: application/json" header on the response.
	// Default of that is "text/plain; charset=utf-8"
	writer.Header().Set("Content-Type", "application/json")

	// Write the JSON to response body.
	writer.Write([]byte(js))
}
