// Package logger is a tiny logging helper built on top of zerolog
// that stores structured fields in a context.Context and attaches them to log events automatically.
package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LOGGER_KEYS is the context key for logger keys
const LOGGER_KEYS = "logger_keys"

// Key is a type for logger context keys
type Key string

// Keys is a set of logger context keys
type Keys map[Key]struct{}

// Logger is the package global logger instance
var Logger = &logger{zlog: log.Logger}

type logger struct {
	zlog zerolog.Logger
}

func updateEventFromContext(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if loggerKeys, ok := ctx.Value(LOGGER_KEYS).(Keys); ok {
		for k := range loggerKeys {
			event.Interface(string(k), ctx.Value(k))
		}
	}

	return event
}

// Debug creates a debug level log event with context fields
//
// You must call Msg/Send on the returned event in order to send the event.
func (this *logger) Debug(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Debug())
}

// Info creates an info level log event with context fields
//
// You must call Msg/Send on the returned event in order to send the event.
func (this *logger) Info(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Info())
}

// Warn creates a warn level log event with context fields
//
// You must call Msg/Send on the returned event in order to send the event.
func (this *logger) Warn(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Warn())
}

// Error creates an error level log event with context fields
//
// You must call Msg/Send on the returned event in order to send the event.
func (this *logger) Error(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Error())
}

// Fatal creates a fatal level log event with context fields.
//
// The os.Exit(1) function is called by the Msg method, which terminates the program immediately.
//
// You must call Msg/Send on the returned event in order to send the event.
func (this *logger) Fatal(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Fatal())
}

func newLogger(zlog zerolog.Logger) *logger {
	return &logger{zlog: zlog}
}

// Init initializes the package logger with the specified log level.
// use global instance of zerolog.Logger from "github.com/rs/zerolog/log" .
func Init(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.DebugLevel

		log.Error().Err(err).Msg("cannot parse log level")
	}

	Logger = newLogger(log.Logger.Level(lvl))
}

// InitWithZerologLogger initializes the package logger with a custom zerolog.Logger instance.
func InitWithZerologLogger(zlog zerolog.Logger) {
	Logger = newLogger(zlog)
}
