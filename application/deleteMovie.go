package application

import (
	"Meow/internal/data"
	"errors"
	"net/http"
)

func (app *Application) deleteMovieHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the movie id from URL.
	id, err := app.readIdParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	// Try delete movie with this id from database with Delete()
	// function. if returned an error check it and send appropriate
	// error response.
	err = app.Models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	// Ok. return 200 status code and success message.
	err = app.writeJSON(writer, http.StatusOK, envelope{"message": "movie deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
