package application

import (
	"context"
	"fmt"

	"user-service/internal/core/domain"
	"user-service/internal/core/dto"
	corerrors "user-service/internal/core/corerrors"
	"user-service/internal/core/ports/outbound"
)

type RegisterUseCase struct {
	repo        outbound.RegisterUserRepository
	hasher      outbound.PasswordHasher
	idGenerator outbound.IDGenerator
}

func NewRegisterUseCase(repo outbound.RegisterUserRepository, hash outbound.PasswordHasher, id outbound.IDGenerator) RegisterUseCase {
	return RegisterUseCase{
		repo:        repo,
		hasher:      hash,
		idGenerator: id,
	}
}

func (uc RegisterUseCase) Register(ctx context.Context, input dto.RegisterInput) (dto.RegisterOutput, error) {

	exist, err := uc.repo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return dto.RegisterOutput{}, fmt.Errorf("check user existence by email: %w", err)
	}
	if exist {
		return dto.RegisterOutput{}, corerrors.ErrUserAlreadyExists
	}

	passwordHashStr, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return dto.RegisterOutput{}, fmt.Errorf("hash password: %w", err)
	}

	id := uc.idGenerator.NewID()
	user := domain.User{
		ID:           id,
		Email:        input.Email,
		PasswordHash: passwordHashStr,
	}

	err = uc.repo.Save(ctx, user)
	if err != nil {
		return dto.RegisterOutput{}, fmt.Errorf("save user: %w", err)
	}

	return dto.RegisterOutput{UserID: id}, nil
}
