package application

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *Application) Serve() error {
	// Declare a http server with settings.
	srv := &http.Server{
		Addr:         app.Config.GetSport(),
		Handler:      app.Routes(),
		ErrorLog:     log.New(app.Logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	// Start a background goroutine.
	go func() {
		// Create a channel for os.Signal value.
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel. Any other signals will not be caught by
		// signal.Notify() and will retain their default behavior.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		// Read the signal from the quit channel. This code will block until a signal is
		// received.
		s := <-quit

		// Log a message to say that the signal has been caught. Notice that we also
		// call the String() method on the signal to get the signal name and include it
		// in the log entry properties.
		app.Logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

		// Exit the application with a 0 (success) status code.
		os.Exit(0)
	}()

	// Logging start serving message.
	app.Logger.PrintInfo("start server", map[string]string{
		"addr": srv.Addr,
		"env":  app.Config.Env,
	})

	// Return error if exist.
	return srv.ListenAndServe()
}
