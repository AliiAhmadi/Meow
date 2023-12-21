package application

import (
	"log"
	"net/http"
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

	// Logging start serving message.
	app.Logger.PrintInfo("start server", map[string]string{
		"addr": srv.Addr,
		"env":  app.Config.Env,
	})

	// Return error if exist.
	return srv.ListenAndServe()
}
