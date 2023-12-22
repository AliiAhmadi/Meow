package application

import (
	"context"
	"errors"
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

	// Create a shutdown error channel for any error returned from
	// Shutdown() function.
	shutdownErr := make(chan error)

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
		app.Logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// Create a context with 5-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}

		app.Logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		app.Wg.Wait()
		shutdownErr <- nil
	}()

	// Logging start serving message.
	app.Logger.PrintInfo("start server", map[string]string{
		"addr": srv.Addr,
		"env":  app.Config.Env,
	})

	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise we will wait to get return value from Shutdown().
	err = <-shutdownErr
	if err != nil {
		return err
	}

	// At this point we know that the graceful shutdown completed successfully and we
	// log a "stopped server" message.
	app.Logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
