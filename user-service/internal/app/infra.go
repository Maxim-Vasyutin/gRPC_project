package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// newPool открывает пул соединений с Postgres и сразу проверяет, что база отвечает.
// Закрывает пул вызывающий (Run), потому что он владеет временем жизни сервера.
func newPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
