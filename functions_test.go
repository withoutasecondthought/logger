package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/withoutasecondthought/logger"
)

func getLoggerKeys(t *testing.T, ctx context.Context) logger.Keys {
	t.Helper()

	raw := ctx.Value(logger.LOGGER_KEYS)
	if raw == nil {
		return nil
	}

	keys, ok := raw.(logger.Keys)
	require.Truef(t, ok, "LOGGER_KEYS has unexpected type %T", raw)

	return keys
}

func keysToStrings(keys logger.Keys) []string {
	out := make([]string, 0, len(keys))
	for k := range keys {
		out = append(out, string(k))
	}

	return out
}

func containsKey(keys logger.Keys, key string) bool {
	if keys == nil {
		return false
	}

	_, ok := keys[logger.Key(key)]

	return ok
}

func TestSetLoggerField_AttachesSingleFieldAndTracksKey(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx2 := logger.SetLoggerField(ctx, "request_id", "r-123")

	assert.Equal(t, "r-123", ctx2.Value(logger.Key("request_id")))

	loggerKeys := getLoggerKeys(t, ctx2)
	require.Len(t, loggerKeys, 1, "expected 1 logger key, got %v", keysToStrings(loggerKeys))
	assert.True(t, containsKey(loggerKeys, "request_id"))
}

func TestSetLoggerFields_AttachesMultipleFieldsAndTracksKeys(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx2 := logger.SetLoggerFields(ctx, map[string]any{
		"user_id": 42,
		"role":    "admin",
	})

	assert.Equal(t, 42, ctx2.Value(logger.Key("user_id")))
	assert.Equal(t, "admin", ctx2.Value(logger.Key("role")))

	loggerKeys := getLoggerKeys(t, ctx2)
	require.Len(t, loggerKeys, 2, "expected 2 logger keys, got %v", keysToStrings(loggerKeys))
	assert.True(t, containsKey(loggerKeys, "user_id"))
	assert.True(t, containsKey(loggerKeys, "role"))
}

func TestSetLoggerField_AccumulatesKeysOverMultipleCalls(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx1 := logger.SetLoggerField(ctx, "request_id", "r-123")
	ctx2 := logger.SetLoggerField(ctx1, "user_id", 42)

	assert.Equal(t, "r-123", ctx2.Value(logger.Key("request_id")))
	assert.Equal(t, 42, ctx2.Value(logger.Key("user_id")))

	loggerKeys := getLoggerKeys(t, ctx2)
	require.Len(t, loggerKeys, 2, "expected 2 logger keys, got %v", keysToStrings(loggerKeys))
	assert.True(t, containsKey(loggerKeys, "request_id"))
	assert.True(t, containsKey(loggerKeys, "user_id"))
}

func TestSetLoggerFields_AccumulatesWithExistingKeys(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx1 := logger.SetLoggerField(ctx, "request_id", "r-123")
	ctx2 := logger.SetLoggerFields(ctx1, map[string]any{
		"user_id": 42,
		"role":    "admin",
	})

	loggerKeys := getLoggerKeys(t, ctx2)
	require.Len(t, loggerKeys, 3, "expected 3 logger keys, got %v", keysToStrings(loggerKeys))
	assert.True(t, containsKey(loggerKeys, "request_id"))
	assert.True(t, containsKey(loggerKeys, "user_id"))
	assert.True(t, containsKey(loggerKeys, "role"))
}

func TestSetPackageWrapper(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx2 := logger.SetPackage(ctx, "example/main")

	assert.Equal(t, "example/main", ctx2.Value(logger.Key("package")))

	loggerKeys := getLoggerKeys(t, ctx2)
	require.Len(t, loggerKeys, 1)
	assert.True(t, containsKey(loggerKeys, "package"))
}

func TestSetFunctionWrapper(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx2 := logger.SetFunction(ctx, "main")

	assert.Equal(t, "main", ctx2.Value(logger.Key("function")))

	loggerKeys := getLoggerKeys(t, ctx2)
	require.Len(t, loggerKeys, 1)
	assert.True(t, containsKey(loggerKeys, "function"))
}

func TestSetLoggerFields_WithEmptyMapPreservesExistingKeys(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx1 := logger.SetLoggerField(ctx, "request_id", "r-123")
	keys1 := getLoggerKeys(t, ctx1)

	ctx2 := logger.SetLoggerFields(ctx1, map[string]any{})
	keys2 := getLoggerKeys(t, ctx2)

	require.Len(t, keys1, len(keys2))
	for k := range keys1 {
		assert.Truef(t, containsKey(keys2, string(k)), "expected key %q to be preserved", k)
	}
}

func TestSetLoggerField_AllowsNilContext(t *testing.T) {
	t.Parallel()

	var ctx context.Context // nil

	ctx2 := logger.SetLoggerField(ctx, "k", "v")

	require.NotNil(t, ctx2)
	assert.Equal(t, "v", ctx2.Value(logger.Key("k")))

	loggerKeys := getLoggerKeys(t, ctx2)
	assert.True(t, containsKey(loggerKeys, "k"))
}

func TestSetLoggerFields_AllowsNilContext(t *testing.T) {
	t.Parallel()

	var ctx context.Context // nil

	ctx2 := logger.SetLoggerFields(ctx, map[string]any{"k": "v"})

	require.NotNil(t, ctx2)
	assert.Equal(t, "v", ctx2.Value(logger.Key("k")))

	loggerKeys := getLoggerKeys(t, ctx2)
	assert.True(t, containsKey(loggerKeys, "k"))
}
