package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger_initLogger(t *testing.T) {
	defer func() {
		require.Nil(t, tearDown())
	}()

	ctx := context.Background()

	t.Run("must_return_logger_with_info_level", func(t *testing.T) {
		logger := initLogger()

		handler := logger.Handler()

		require.False(t, handler.Enabled(ctx, slog.LevelDebug))
		require.True(t, handler.Enabled(ctx, slog.LevelInfo))
	})

	t.Run("must_return_logger_with_debug_level", func(t *testing.T) {
		err := os.Setenv(envDebug, "true")
		require.Nil(t, err)

		logger := initLogger()

		handler := logger.Handler()

		require.True(t, handler.Enabled(ctx, slog.LevelDebug))
	})

}

func TestLogger_isDebug(t *testing.T) {
	defer func() {
		require.Nil(t, tearDown())
	}()

	require.False(t, isDebug())

	err := os.Setenv(envDebug, "true")
	require.Nil(t, err)
	require.True(t, isDebug())

	err = os.Setenv(envDebug, "1")
	require.Nil(t, err)
	require.True(t, isDebug())

	err = os.Setenv(envDebug, "false")
	require.Nil(t, err)
	require.False(t, isDebug())

	err = os.Setenv(envDebug, "0")
	require.Nil(t, err)
	require.False(t, isDebug())

}

func tearDown() error {
	return os.Unsetenv(envDebug)
}
