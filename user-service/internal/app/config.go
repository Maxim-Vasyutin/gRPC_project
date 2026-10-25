package app

import (
	"log"
	"os"
	"time"
)

// Значения по умолчанию — только для локальной разработки.
const (
	defaultDSN      = "postgres://user:password@localhost:5434/users?sslmode=disable"
	defaultGRPCPort = ":50051"
	devJWTSecret    = "dev-secret-do-not-use-in-prod"
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
	)

type config struct {
	dsn        string
	grpcPort   string
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// loadConfig читает настройки из переменных окружения.
// Если переменная не задана — берётся значение для локальной разработки.
func loadConfig() config {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET не задан, используется dev-секрет. Не для прода.")
		secret = devJWTSecret
	}

	return config{
		dsn:        getEnv("DB_DSN", defaultDSN),
		grpcPort:   getEnv("GRPC_PORT", defaultGRPCPort),
		jwtSecret:  []byte(secret),
		accessTTL:  accessTokenTTL,
		refreshTTL: refreshTokenTTL,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
