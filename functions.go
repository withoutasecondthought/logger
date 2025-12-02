package logger

import "context"

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

func SetPackage(ctx context.Context, packageName string) context.Context {
	return SetLoggerField(ctx, "package", packageName)
}

func SetFunction(ctx context.Context, functionName string) context.Context {
	return SetLoggerField(ctx, "function", functionName)
}
