// logger.go
package utils

import (
	"fmt"
	"io"
	"log/slog"
)

type LogLevel slog.Level

var LogLevelMap = map[string]slog.Level{
	"ERROR":   slog.LevelError,
	"WARNING": slog.LevelWarn,
	"INFO":    slog.LevelInfo,
	"DEBUG":   slog.LevelDebug,
}

// Logger is an abstraction for logging.
type Logger interface {
	Error(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// Log is an implementation of a Logger using slog.
type Log struct {
	logger *slog.Logger
}

// NewSlogLogger creates a new SlogLogger instance.
func NewLogger(h slog.Handler) *Log {
	return &Log{
		logger: slog.New(h),
	}
}

// NewWriter creates a new slog handler, text or json
func NewLogHandler(handlerType string, handlerWriter io.Writer, handlerOptions *slog.HandlerOptions) (slog.Handler, error) {
	switch handlerType {
	case "text":
		return slog.NewTextHandler(handlerWriter, handlerOptions), nil
	case "json":
		return slog.NewJSONHandler(handlerWriter, handlerOptions), nil
	default:
		return nil, fmt.Errorf("failed to create logging handler. Possible values: text, json")
	}
}

// Err logs a message at the ERR level.
func (l *Log) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

// Warning logs a message at the WARNING level.
func (l *Log) Warning(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

// Info logs a message at the INFO level.
func (l *Log) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

// Debug logs a message at the DEBUG level.
func (l *Log) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}
