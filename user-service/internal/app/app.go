// Package app — composition root: собирает все слои сервиса и управляет его жизнью.
//
// Файлы пакета повторяют фазы запуска:
//
//	config.go — откуда берутся настройки
//	infra.go  — подключение к внешнему миру (БД)
//	deps.go   — обёртки над конструкторами для fx: достают поля из config
//	grpc.go   — транспорт: gRPC-сервер, интерсепторы, хендлеры
package app

import (
	"user-service/internal/adapters/outbound/postgres"
	"user-service/internal/application"
	"user-service/internal/core/ports/outbound"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Run — оглавление запуска.  fx сам устанавливает связи и порядок
func Run(logger *zap.Logger) {
	fx.New(
		fx.Supply(logger),
		fx.Provide(
			loadConfig,
			newPool,
			newGRPCServer,
			fx.Annotate(
				postgres.NewPostgresUserRepository,
				fx.As(new(outbound.RegisterUserRepository)),
				fx.As(new(outbound.LoginUserRepository)),
			),
			fx.Annotate(
				postgres.NewPostgresSessionRepository,
				fx.As(
					new(outbound.SessionRepository),
				),
			),
			fx.Annotate(
				newPasswordHasher,
				fx.As(
					new(outbound.PasswordHasher),
				),
			),
			fx.Annotate(
				newTokenIssuer,
				fx.As(
					new(outbound.TokenIssuer),
				),
			),
			fx.Annotate(
				newIDGenerator,
				fx.As(
					new(outbound.IDGenerator),
				),
			),
			newSessionIssuer,
			application.NewRegisterUseCase,
			application.NewLoginUseCase,
			newRefreshUseCase,
			application.NewLogoutUseCase,
		),
		fx.Invoke(
			registerServer,
		),
	).Run()
}

/*
Про fx
Сам соблюдает порядок

Provide() - список функций без порядка. Ф-я не вызывается, пока её результат не потребуется дальше.
Supply() - кладёт в граф уже готовое значение
Invoke() - при старте вызывает этот блок обязательно. Тут и возникает цепочка: что-то надо при запуске - вызывается Provide()
Lifecycle - параметр, показывает жизненный цикл fx. внутри есть параметры OnStart и OnStop
New().Run() - Собрать и вызывать (и потом откатиться)

*/
