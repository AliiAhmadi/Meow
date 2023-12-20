package jlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Define a Level type to represent the severity level for a log entry.
type Level int8

// Initialize constants which represent a specific severity level.
const (
	LevelInfo  Level = iota // 0
	LevelError              // 1
	LevelFatal              // 2
	LevelOff                // 3
)

// Return a human-friendly string for the severity level.
func (level Level) String() string {
	switch level {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Define a custom Logger type. This holds the output destination that the log entries
// will be written to, the minimum severity level that log entries will be written for,
// plus a mutex for coordinating the writes.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// Return a new Logger instance which writes log entries at or above a minimum severity
// level to a specific output destination.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// Declare PrintInfo helper methods for writing log entries at the INFO level.
func (logger *Logger) PrintInfo(message string, properties map[string]string) {
	logger.print(LevelInfo, message, properties)
}

// Declare PrintError helper methods for writing log entries at the ERROR level.
func (logger *Logger) PrintError(err error, properties map[string]string) {
	logger.print(LevelError, err.Error(), properties)
}

// Declare PrintFatal helper methods for writing log entries at the FATAL level.
func (logger *Logger) PrintFatal(err error, properties map[string]string) {
	logger.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

// Print is an internal method for writing the log entry.
func (logger *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	// If the severity level of the log entry is below the minimum severity for the
	// logger, then return with no further action.
	if logger.minLevel > level {
		return 0, nil
	}

	// Declare an struct holding the data for the log entry.
	st := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	// Include a stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		st.Trace = string(debug.Stack())
	}

	// Marshal the anonymous struct to JSON and store it in the line variable.
	line, err := json.Marshal(st)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	// Lock the mutex so that no two writes to the output destination can happen
	// concurrently.
	logger.mu.Lock()
	defer logger.mu.Unlock()

	// Log the write entry.
	return logger.out.Write(append(line, '\n'))
}

// We also implement a Write() method on our Logger type so that it satisfies the
// io.Writer interface.
func (logger *Logger) Write(message []byte) (int, error) {
	return logger.print(LevelError, string(message), nil)
}
