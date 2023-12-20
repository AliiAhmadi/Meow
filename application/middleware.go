package application

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// recoverPanic() middleware
func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response.
				writer.Header().Set("Connection", "close")

				// The value returned by recover() has the type interface{}, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				app.serverErrorResponse(writer, request, fmt.Errorf("%s", err))
			}
		}()

		// Call next
		next.ServeHTTP(writer, request)
	})
}

// rateLimit() middleware
func (app *Application) rateLimit(next http.Handler) http.Handler {
	// Declare a mutex and a map to hold the clients' IP addresses and rate limiters.
	var (
		mu      sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	// The function we are returning is a closure, which 'closes over' the limiter
	// variable.
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Extract the client's IP address from the request.
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(writer, request, err)
			return
		}

		// Lock the mutex to prevent this code from being executed concurrently.
		mu.Lock()

		// Check to see if the IP address already exists in the map. If it doesn't, then
		// initialize a new rate limiter and add the IP address and limiter to the map.
		if _, found := clients[ip]; !found {
			clients[ip] = rate.NewLimiter(2, 4)
		}

		// Call the Allow() method on the rate limiter for the current IP address. If
		// the request isn't allowed, unlock the mutex and send a 429 Too Many Requests
		// response, just like before.
		if !clients[ip].Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(writer, request)
			return
		}

		// Very importantly, unlock the mutex before calling the next handler in the
		// chain.
		mu.Unlock()

		// Call next.
		next.ServeHTTP(writer, request)
	})
}
