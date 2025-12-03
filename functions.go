package logger

import "context"

// SetLoggerField adds a key-value pair to the logger context.
func SetLoggerField(ctx context.Context, key string, value any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	loggerKeys, ok := ctx.Value(LOGGER_KEYS).(Keys)

	newLoggerKeys := make(Keys)
	if ok {
		for k := range loggerKeys {
			newLoggerKeys[k] = struct{}{}
		}
	}

	newLoggerKeys[Key(key)] = struct{}{}

	return context.WithValue(context.WithValue(
		ctx, Key(key), value,
	), LOGGER_KEYS, newLoggerKeys)
}

// SetLoggerFields adds multiple key-value pairs to the logger context.
func SetLoggerFields(ctx context.Context, fields map[string]any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	loggerKeys, ok := ctx.Value(LOGGER_KEYS).(Keys)

	newLoggerKeys := make(Keys)
	if ok {
		for k := range loggerKeys {
			newLoggerKeys[k] = struct{}{}
		}
	}

	for k, v := range fields {
		newLoggerKeys[Key(k)] = struct{}{}
		//nolint:fatcontext
		ctx = context.WithValue(ctx, Key(k), v)
	}

	return context.WithValue(ctx, LOGGER_KEYS, newLoggerKeys)
}

// SetPackage adds the package name to the logger context.
//
// Shorthand for SetLoggerField with key "package".
func SetPackage(ctx context.Context, packageName string) context.Context {
	return SetLoggerField(ctx, "package", packageName)
}

// SetFunction adds the function name to the logger context.
//
// Shorthand for SetLoggerField with key "function".
func SetFunction(ctx context.Context, functionName string) context.Context {
	return SetLoggerField(ctx, "function", functionName)
}

// SetPackageAndFunction adds both package and function names to the logger context.
//
// Shorthand for SetLoggerFields with keys "package" and "function".
func SetPackageAndFunction(ctx context.Context, packageName, functionName string) context.Context {
	return SetLoggerFields(ctx, map[string]any{"package": packageName, "function": functionName})
}
