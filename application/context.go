package application

import (
	"Meow/internal/data"
	"context"
	"net/http"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

const userContextKey = contextKey("user")

// The contextSetUser() method returns a new copy of the request with the provided
// User struct added to the context.
func (app *Application) contextSetUser(request *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(request.Context(), userContextKey, user)
	return request.WithContext(ctx)
}

func (app *Application) contextGetUser(request *http.Request) *data.User {
	user, ok := request.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
