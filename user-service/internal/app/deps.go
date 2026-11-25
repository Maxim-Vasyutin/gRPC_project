package app

import (
	"user-service/internal/adapters/outbound/hasher"
	"user-service/internal/adapters/outbound/idgen"
	"user-service/internal/adapters/outbound/token"
	"user-service/internal/application"
	"user-service/internal/core/ports/outbound"
)

func newPasswordHasher(cfg config) (hasher.BcryptHasher, error) {
	return hasher.NewBcryptHasher(cfg.bcryptCost)
}

func newTokenIssuer(cfg config) *token.JWTIssuer {
	return token.NewJWTIssuer(cfg.jwtSecret, cfg.accessTTL)
}

func newIDGenerator() idgen.UUIDGenerator {
	return idgen.UUIDGenerator{}
}

// use cases
func newSessionIssuer(tokens outbound.TokenIssuer, sessionRepo outbound.SessionRepository, cfg config) *application.SessionIssuer {
	return application.NewSessionIssuer(tokens, sessionRepo, cfg.refreshTTL)
}

func newRefreshUseCase(session outbound.SessionRepository, tokenIssuer outbound.TokenIssuer, cfg config) *application.RefreshUseCase {
	return application.NewRefreshUseCase(session, tokenIssuer, cfg.refreshTTL)
}
