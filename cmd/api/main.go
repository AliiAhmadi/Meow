package main

import (
	"Meow/application"
	"Meow/config"
	"Meow/internal/data"
	jlog "Meow/log"
	"flag"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
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

	// Read the DSN from command line flags into config struct.
	flag.StringVar(&cfg.DB.DSN, "dsn", os.Getenv("dsn"), "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags.
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := jlog.New(os.Stdout, jlog.LevelInfo)

	// Get connection pool from OpenDB() function.
	// If any error exists, we should exit application immediately.
	db, err := config.OpenDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	// defering closing database pool.
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.PrintInfo("database connection pool established", nil)

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	application := &application.Application{
		Config:  cfg,
		Logger:  logger,
		Version: version,
		Models:  data.NewModels(db),
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
	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.Env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}
