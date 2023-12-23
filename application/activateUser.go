package application

import (
	"Meow/internal/data"
	"Meow/internal/validator"
	"errors"
	"net/http"
)

func (app *Application) activateUserHandler(writer http.ResponseWriter, request *http.Request) {
	// Get activation code from user request.
	var input struct {
		TokenPlain string `json:"token"`
	}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	// Validate the plaintext token provided by the client.
	v := validator.New()

	if validator.ValidateTokenPlaintext(v, input.TokenPlain); !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	user, err := app.Models.Users.GetForToken(data.ScopeActivation, input.TokenPlain)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(writer, request, v.Errors)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.Models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.Models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(writer, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
}
