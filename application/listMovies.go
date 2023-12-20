package application

import (
	"Meow/internal/data"
	"Meow/internal/validator"
	"net/http"
)

// Get a list of all movies based on query parameters provided by user.
func (app *Application) listMoviesHandler(writer http.ResponseWriter, request *http.Request) {
	// To keep things consistent with our other handlers, we'll define an input struct
	// to hold the expected values from the request query string.
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	// Initialize a new validator instance.
	v := validator.New()

	// r.URL.Query() will return url.Values (a map containing the query string data)
	qs := request.URL.Query()

	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 10, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "id")

	// Add safe sorts to safeSort slice for validation in follow.
	input.Filters.SortSafeList = []string{
		"id",
		"-id",
		"title",
		"-title",
		"year",
		"-year",
		"runtime",
		"-runtime",
	}

	// Validating filters and if any error exists
	// write that in v.
	validator.ValidateFilters(v, input.Filters)

	// Check for validation errors. if any error
	// exist send failedValidationResponse() response.
	if !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	// Get Movies based on query parameters by calling GetAll() method.
	movies, err := app.Models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"movies": movies}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
}
