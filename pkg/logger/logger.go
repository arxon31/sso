package logger

import (
	"log/slog"
	"os"
	"strconv"
)

const envDebug = "DEBUG"

var Logger *slog.Logger

func init() {
	Logger = initLogger()
}

func initLogger() *slog.Logger {
	handlerOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	if isDebug() {
		handlerOpts.Level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, handlerOpts)

	return slog.New(handler)
}

func isDebug() bool {
	env, ok := os.LookupEnv(envDebug)
	if !ok {
		return false
	}

	debug, err := strconv.ParseBool(env)
	if err != nil {
		return false
	}

	return debug
}
