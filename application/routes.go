package application

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	READ_PERMISSION  = "movies:read"
	WRITE_PERMISSION = "movies:write"
)

func (app *Application) Routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Set not found error handler for router.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission(WRITE_PERMISSION, app.createNewMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission(READ_PERMISSION, app.showMovieHandler))
	router.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.requirePermission(WRITE_PERMISSION, app.updateMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission(WRITE_PERMISSION, app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission(WRITE_PERMISSION, app.deleteMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission(READ_PERMISSION, app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Return the http.Handler instance.
	// Wrapped with recoverPanic() middleware.
	// Also wrap that with rateLimit() middleware. (v2)
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
