package app

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Значения по умолчанию — только для локальной разработки.
const (
	defaultDSN        = "postgres://user:password@localhost:5434/users?sslmode=disable"
	defaultGRPCPort   = ":50051"
	devJWTSecret      = "dev-secret-do-not-use-in-prod"
	accessTokenTTL    = 15 * time.Minute
	refreshTokenTTL   = 30 * 24 * time.Hour
	defaultBcryptCost = bcrypt.DefaultCost
)

type config struct {
	dsn        string
	grpcPort   string
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	bcryptCost int
}

// loadConfig читает настройки из переменных окружения.
// Если переменная не задана — берётся значение по умолчанию для локальной разработки.
func loadConfig(logger *zap.Logger) (config, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logger.Warn("JWT_SECRET не задан, используется dev-секрет. Не для прода.")
		secret = devJWTSecret
	}

	cost := defaultBcryptCost
	if v := os.Getenv("BCRYPT_COST"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return config{}, fmt.Errorf("BCRYPT_COST must be an integer: %w", err)
		}
		cost = parsed
	}

	return config{
		dsn:        getEnv("DB_DSN", defaultDSN),
		grpcPort:   getEnv("GRPC_PORT", defaultGRPCPort),
		jwtSecret:  []byte(secret),
		accessTTL:  accessTokenTTL,
		refreshTTL: refreshTokenTTL,
		bcryptCost: cost,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
