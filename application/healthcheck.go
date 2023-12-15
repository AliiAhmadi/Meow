package application

import (
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *Application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {

	// Create a map for data we want to send.
	data := map[string]string{
		"status":      "available",
		"environment": app.Config.Env,
		"version":     app.Version,
	}

	// Using writeJSON() helper to write data and set
	// headers to HTTP response body.
	err := app.writeJSON(writer, http.StatusOK, envelope{"info": data}, nil)
	if err != nil {
		app.Logger.Println(err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}
}
