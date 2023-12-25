package main

import (
	"Meow/application"
	"Meow/config"
	"Meow/internal/data"
	jlog "Meow/log"
	"Meow/mailer"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Declare a string containing the application version number.
// Version number hard-coded constant.

var (
	buildTime string
	version   = "1.0.0"
)

func expvarValues(db *sql.DB) {
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))
}

func main() {
	// Declare an instance of config struct.
	cfg := new(config.Config)

	// expvar customization

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
	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Rate limiter enable-disable")

	flag.StringVar(&cfg.Smtp.Host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", "-", "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", "-", "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", "Meow <no-reply@meow.com>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(value string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(value)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Dispaly version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

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

	expvarValues(db)
	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	application := &application.Application{
		Config:  cfg,
		Logger:  logger,
		Version: version,
		Models:  data.NewModels(db),
		Mailer: mailer.New(
			cfg.Smtp.Host,
			cfg.Smtp.Port,
			cfg.Smtp.Username,
			cfg.Smtp.Password,
			cfg.Smtp.Sender,
		),
	}

	// Start server with serve() method in app instance.
	err = application.Serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
