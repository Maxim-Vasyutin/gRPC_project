package core

import (
	"context"
)

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user User) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}
