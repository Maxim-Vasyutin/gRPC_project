package main

import (
	"fmt"
	"os"

	"user-service/internal/app"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Логгер создаётся один раз здесь и дальше передаётся явно.
	logger, err := newLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "create logger: %v\n", err)
		os.Exit(1)
	}

	app.Run(logger)
	_ = logger.Sync()
}

func newLogger() (*zap.Logger, error) {
	if os.Getenv("APP_ENV") == "production" {
		return zap.NewProduction() // JSON для прода
	}
	cfg := zap.NewDevelopmentConfig() // читаемый формат для разработки
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
}
