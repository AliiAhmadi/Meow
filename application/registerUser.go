package application

import (
	"Meow/internal/data"
	"Meow/internal/validator"
	"errors"
	"net/http"
	"time"
)

func (app *Application) registerUserHandler(writer http.ResponseWriter, request *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body to `input` struct.
	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	// Create a new user from input data.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if validator.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	// Insert the user data into the database.
	err = app.Models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "this email address already exists")
			app.failedValidationResponse(writer, request, v.Errors)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.Models.Permissions.AddForUsers(user.ID, READ_PERMISSION)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	// Generate a token for acvivation.
	token, err := app.Models.Tokens.New(user.ID, 1*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	// Run sending email and panic recover for that in background.
	app.background(func() {

		data := map[string]interface{}{
			"userID":          user.ID,
			"activationToken": token.Plaintext,
		}

		err := app.Mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.Logger.PrintError(err, nil)
		}
	})

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(writer, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
}
