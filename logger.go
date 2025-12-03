package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const LOGGER_KEYS = "logger_keys"

type Key string

type Keys map[Key]struct{}

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

func (this *logger) Debug(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Debug())
}

func (this *logger) Info(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Info())
}

func (this *logger) Warn(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Warn())
}

func (this *logger) Error(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Error())
}

func (this *logger) Fatal(ctx context.Context) *zerolog.Event {
	return updateEventFromContext(ctx, this.zlog.Fatal())
}

func newLogger(zlog zerolog.Logger) *logger {
	return &logger{zlog: zlog}
}

func Init(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.DebugLevel

		log.Error().Err(err).Msg("cannot parse log level")
	}

	Logger = newLogger(log.Logger.Level(lvl))
}

func InitWithZerologLogger(zlog zerolog.Logger) {
	Logger = newLogger(zlog)
}
