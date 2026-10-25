package postgres

import (
	"context"
	"errors"
	"user-service/internal/core/corerrors"
	"user-service/internal/core/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) PostgresUserRepository {
	return PostgresUserRepository{
		pool: pool,
	}
}

func (r PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM users
		WHERE email = $1
		)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func (r PostgresUserRepository) Save(ctx context.Context, user domain.User) error {
	query := `
	INSERT INTO users (id, email, password_hash)
	VALUES ($1, $2, $3)
	`

	_, err := r.pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (r PostgresUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
	SELECT id, email, password_hash
	FROM users
	WHERE email = $1
	`

	var user domain.User

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, corerrors.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

