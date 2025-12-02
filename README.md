# Logger

A tiny logging helper built on top of zerolog that stores structured fields in a
`context.Context` and attaches them to log events automatically.

This package provides a global `Logger` instance with helper functions to add
fields into a context and to create log events that include those fields.

Highlights:
- Keep request-scoped fields in `context.Context` and attach them to each log
  event.
- Simple helpers for adding a single field or multiple fields to the context.
- Small, zerolog-based logger with methods: Debug, Info, Warn, Error, Fatal.
- A simple `Init(level string)` function to configure the global log level.

## Installation

Use Go modules as usual:

    go get github.com/withoutaseccondthought/logger

## Quick start

Example `main.go`:

```go
package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/withoutaseccondthought/logger"
)

func main() {
	// Optional: set global zerolog level by string (e.g. "debug", "info").
	logger.Init("debug")

	// Create a context and attach fields that should travel with it.
	ctx := context.Background()
	ctx = logger.SetLoggerField(ctx, "request_id", "r-123")
	ctx = logger.SetPackage(ctx, "example/main")
	ctx = logger.SetFunction(ctx, "main")

	// Log with the package's Logger. Fields from ctx will be attached.
	logger.Logger.Info(ctx).Msg("starting up")

	// Add many fields at once
	ctx = logger.SetLoggerFields(ctx, map[string]any{
		"user_id": 42,
		"role":    "admin",
	})

	logger.Logger.Debug(ctx).Msg("user info attached")

	// Use zerolog chaining as usual — fields from context are added first.
	logger.Logger.Error(ctx).Str("extra", "value").Msg("an error occurred")

	_ = log.Logger // keep the import if you need direct zerolog usage
}
```

## Notes and best practices

- Avoid storing large or long-lived objects in `context.Context`. Use it only
  for request-scoped metadata (IDs, small values) to keep logs concise.
- Keys tracked by the package are stored under the `LOGGER_KEYS` constant.
  The package uses string keys when writing fields to zerolog; use unique
  names to avoid collisions.
- The package delegates formatting and output behavior to zerolog. If you
  prefer human-friendly console output in development, configure zerolog's
  global logger accordingly before calling `logger.Init` or replace the
  underlying `log.Logger`.

## API reference

- func Init(level string)
  - Parse `level` with zerolog's ParseLevel and set the global log level.
  - If parsing fails, the logger falls back to `Debug` level and logs the
    parsing error.

- func SetLoggerField(ctx context.Context, key string, value interface{}) context.Context
  - Returns a new context with `key` set to `value` and tracks the key so the
    logger can copy the value into events.

- func SetLoggerFields(ctx context.Context, fields map[string]any) context.Context
  - Like `SetLoggerField`, but adds multiple key/value pairs from `fields`.

- func SetPackage(ctx context.Context, packageName string) context.Context
  - Convenience wrapper for `SetLoggerField(ctx, "package", packageName)`.

- func SetFunction(ctx context.Context, functionName string) context.Context
  - Convenience wrapper for `SetLoggerField(ctx, "function", functionName)`.

- var Logger *logger
  - Global logger instance. Methods on the logger produce zerolog events and
    automatically attach the tracked context fields:
    - Debug(ctx context.Context) *zerolog.Event
    - Info(ctx context.Context) *zerolog.Event
    - Warn(ctx context.Context) *zerolog.Event
    - Error(ctx context.Context) *zerolog.Event
    - Fatal(ctx context.Context) *zerolog.Event

## Contributing

PRs and issues are welcome. Keep changes small and focused. Add tests for new
behaviour where appropriate.