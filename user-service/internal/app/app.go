// Package app — composition root: собирает все слои сервиса и управляет его жизнью.
//
// Файлы пакета повторяют фазы запуска:
//
//	config.go — откуда берутся настройки
//	infra.go  — подключение к внешнему миру (БД)
//	deps.go   — сборка адаптеров и use cases
//	grpc.go   — транспорт: gRPC-сервер, интерсепторы, хендлеры
package app

import (
	"context"

	"go.uber.org/zap"
)

// Run — оглавление запуска. Порядок строк = порядок фаз.
func Run(logger *zap.Logger) error {
	ctx := context.Background()

	cfg, err := loadConfig(logger)
	if err != nil {
		return err
	}

	pool, err := newPool(ctx, cfg.dsn)
	if err != nil {
		return err
	}
	defer pool.Close() // пул живёт столько же, сколько Run, то есть сколько работает сервер

	uc, err := newUseCases(pool, cfg)
	if err != nil {
		return err
	}

	grpcServer, err := newGRPCServer(uc, logger)
	if err != nil {
		return err
	}

	return serve(grpcServer, cfg.grpcPort, logger)
}
