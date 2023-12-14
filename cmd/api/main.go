package main

import (
	"Meow/application"
	"Meow/config"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

// Declare a string containing the application version number.
// Version number hard-coded constant.
const version = "1.0.0"

func main() {
	// Declare an instance of config struct.
	cfg := new(config.Config)

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "Log: ", log.Ldate|log.Ltime)

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	application := &application.Application{
		Config:  cfg,
		Logger:  logger,
		Version: version,
	}

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler. use the httprouter instance returned by app.routes() as the server handler.
	srv := &http.Server{
		Addr:         cfg.GetSport(),
		Handler:      application.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the http server.
	logger.Printf("starting %s server on %s", cfg.Env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
