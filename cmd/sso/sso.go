package main

import (
	"github.com/arxon31/sso/pkg/logger"
	"log/slog"
	"os"
)

var Build string

func main() {
	logger.Logger.Info("starting app", slog.String("build", Build))
	err := run()
	if err != nil {
		logger.Logger.Error("app exited with error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Logger.Info("app exited without errors")
}

func run() error {
	return nil
}
