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
			newUseCases,
			newGRPCServer,
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
Supply() - показывает уже готовое значение
Invoke() - при старте вызывает этот блок обязательно. Тут и возникает цепочка: что-то надо при запуске - вызывается Provide()
Lifecycle() - показывает жизненный цикл fx. внутри есть параметры OnStart и OnStop
New().Run() - Собрать и вызывать (и потом откатиться)

*/
