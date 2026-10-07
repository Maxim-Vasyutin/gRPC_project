package application

import (
	"context"
	"fmt"

	"gRPC_project/user-service/internal/application/dto"
	"gRPC_project/user-service/internal/application/ports/outbound"
	"gRPC_project/user-service/internal/core"
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
		return dto.RegisterOutput{}, ErrUserAlreadyExists
	}

	passwordHashStr, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return dto.RegisterOutput{}, fmt.Errorf("hash password: %w", err)
	}

	id := uc.idGenerator.NewID()
	user := core.User{
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
