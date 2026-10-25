package postgres

import (
	"context"
	"errors"
	"user-service/internal/core/corerrors"
	"user-service/internal/core/domain"
	"user-service/internal/core/ports/outbound"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSessionRepository struct {
	pool *pgxpool.Pool
}

var _ outbound.SessionRepository = PostgresSessionRepository{}


func NewPostgresSessionRepository(pool *pgxpool.Pool) PostgresSessionRepository {
	return PostgresSessionRepository{
		pool: pool,
	}
}

func (r PostgresSessionRepository) Store(ctx context.Context, session domain.Session) error {
	query := `
	INSERT INTO sessions (refresh_token, user_id, expires_at)
	VALUES ($1, $2, $3)
	`

	_, err := r.pool.Exec(ctx, query, session.RefreshToken, session.UserID, session.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresSessionRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error) {
	query := `
	SELECT user_id, refresh_token, expires_at
	FROM sessions
	WHERE refresh_token = $1 AND expires_at > now()
	`

	var session domain.Session

	err := r.pool.QueryRow(
		ctx, query, refreshToken,
	).Scan(
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Session{}, corerrors.ErrSessionNotFound
	}

	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r PostgresSessionRepository) Replace(ctx context.Context, oldRefreshToken string, newSession domain.Session) error {
	return errors.New("not implemented")
}

func (r PostgresSessionRepository) RevokeByRefreshToken(ctx context.Context, refreshToken string) error {
	return errors.New("not implemented")
}
