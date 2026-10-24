package app

import (
	"user-service/internal/adapters/outbound/hasher"
	"user-service/internal/adapters/outbound/idgen"
	"user-service/internal/adapters/outbound/memory"
	"user-service/internal/adapters/outbound/postgres"
	"user-service/internal/adapters/outbound/token"
	"user-service/internal/application"
	"user-service/internal/core/ports/inbound"

	"github.com/jackc/pgx/v5/pgxpool"
)

// useCases — все сценарии сервиса, собранные и готовые к работе.
// Поля имеют тип inbound-портов: транспорту не важно, какая реализация внутри.
type useCases struct {
	register inbound.Register
	login    inbound.Login
	// refresh, logout — когда появятся хендлеры
}

// newUseCases создаёт outbound-адаптеры и вставляет их в use cases.
// Это единственное место, где конкретные типы (Postgres, bcrypt, JWT)
// встречаются с портами ядра.
func newUseCases(pool *pgxpool.Pool, cfg config) useCases {
	// ── Outbound-адаптеры ───────────────────────────────
	userRepo := postgres.NewPostgresUserRepository(pool)
	pwHasher := hasher.BcryptHasher{}
	idGen := idgen.UUIDGenerator{}
	issuer := token.NewJWTIssuer(cfg.jwtSecret, cfg.accessTTL)
	sessionRepo := memory.NewSessionRepository() // ВРЕМЕННО: заменить на Postgres-версию

	// ── Use cases ───────────────────────────────────────
	registerUC := application.NewRegisterUseCase(userRepo, pwHasher, idGen)
	loginUC := application.NewLoginUseCase(userRepo, pwHasher, issuer, sessionRepo)

	return useCases{
		register: registerUC,
		login:    &loginUC, // & — потому что метод Login объявлен на *LoginUseCase
	}
}
