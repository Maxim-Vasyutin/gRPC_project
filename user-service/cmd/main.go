package main

import (
	"fmt"
	"os"

	"user-service/internal/app"

	"go.uber.org/zap"
)

func main() {
	// Логгер создаётся один раз здесь и дальше передаётся явно.
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "create logger: %v\n", err)
		os.Exit(1)
	}

	runErr := app.Run(logger)
	if runErr != nil {
		logger.Error("service stopped with error", zap.Error(runErr))
	}
	_ = logger.Sync() // os.Exit не выполняет defer, поэтому Sync вызываем явно
	if runErr != nil {
		os.Exit(1)
	}
}
