package core

import (
	"context"
	"errors"
	"fmt"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID string
}

type RegisterUseCase struct {
	repo        UserRepository
	hasher      PasswordHasher
	idGenerator IDGenerator
}

func NewRegisterUseCase(repo UserRepository, hash PasswordHasher, id IDGenerator) RegisterUseCase {
	return RegisterUseCase{
		repo:        repo,
		hasher:      hash,
		idGenerator: id,
	}
}

func (uc RegisterUseCase) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {

	exist, err := uc.repo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("check user existence by email: %w", err)
	}
	if exist {
		return RegisterOutput{}, ErrUserAlreadyExists
	}

	passwordHashStr, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("hash password: %w", err)
	}

	id := uc.idGenerator.NewID()
	user := User{
		ID:           id,
		Email:        input.Email,
		PasswordHash: passwordHashStr,
	}

	err = uc.repo.Save(ctx, user)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("save user: %w", err)
	}

	return RegisterOutput{UserID: id}, nil
}
