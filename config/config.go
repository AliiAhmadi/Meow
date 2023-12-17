package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

func (cfg *Config) GetSport() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

// The OpenDB() function returns a sql.DB connection pool.
func OpenDB(cfg *Config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool.
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections in the pool.
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	// Set the maximum number of idle connections in the pool.
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	// Using time.ParseDuration() function to convert idle timeout duration
	// to time.Duration type.
	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set the maximum idle time in the pool.
	db.SetConnMaxIdleTime(duration)

	// Create a context with 5 second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return connection pool.
	return db, nil
}
