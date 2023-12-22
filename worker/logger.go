package worker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

// Printf implements internal.Logging.
func (*Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	panic("unimplemented")
}

func NewLogger() *Logger {
	return &Logger{}
}

// print the logs
func (logger *Logger) Print(level zerolog.Level, args ...interface{}) {
	// merge all objects interface into one string by fmt.Sprint
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

// implement logs in the Logger interface
// Debug logs a message at Debug level.
func (logger *Logger) Debug(args ...interface{}) {
	logger.Print(zerolog.DebugLevel, args...)
}

// Info logs a message at Info level.
func (logger *Logger) Info(args ...interface{}) {
	logger.Print(zerolog.InfoLevel, args...)
}

// Warn logs a message at Warning level.
func (logger *Logger) Warn(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}

// Error logs a message at Error level.
func (logger *Logger) Error(args ...interface{}) {
	logger.Print(zerolog.ErrorLevel, args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (logger *Logger) Fatal(args ...interface{}) {
	logger.Print(zerolog.FatalLevel, args...)
}
