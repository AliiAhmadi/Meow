package application

import (
	"Meow/internal/data"
	"Meow/internal/validator"
	"errors"
	"net/http"
	"time"
)

func (app *Application) createAuthenticationTokenHandler(writer http.ResponseWriter, request *http.Request) {
	// Parse the `email` and `password` from request body.
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	// Validation for input fields
	v := validator.New()

	validator.ValidateEmail(v, input.Email)
	validator.ValidatePasswordPlainText(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	// Search for user in database
	user, err := app.Models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(writer, request)
		return
	}

	token, err := app.Models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	err = app.writeJSON(writer, http.StatusCreated, envelope{"auth_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
}
