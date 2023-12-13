package main

import (
	"fmt"
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *Application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "status: available")
	fmt.Fprintf(writer, "environment: %s\n", app.Config.Env)
	fmt.Fprintf(writer, "version: %s\n", version)
}
