package application

import (
	"context"
	"fmt"

	corerrors "user-service/internal/core/corerrors"
	"user-service/internal/core/domain"
	"user-service/internal/core/dto"
	"user-service/internal/core/ports/outbound"
)

type RegisterUseCase struct {
	repo        outbound.RegisterUserRepository
	hasher      outbound.PasswordHasher
	idGenerator outbound.IDGenerator
	session     *SessionIssuer
}

func NewRegisterUseCase(
	repo outbound.RegisterUserRepository,
	hash outbound.PasswordHasher,
	id outbound.IDGenerator,
	session *SessionIssuer) *RegisterUseCase {
	return &RegisterUseCase{
		repo:        repo,
		hasher:      hash,
		idGenerator: id,
		session:     session,
	}
}

func (uc *RegisterUseCase) Register(ctx context.Context, input dto.RegisterInput) (dto.RegisterOutput, error) {

	// В коре написал метод приведения к одному регистру для почты
	emailLower := domain.NormalizeEmail(input.Email)
	loginLower := domain.NormalizeLogin(input.Login)

	exist, err := uc.repo.ExistsByEmail(ctx, emailLower)
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
		Email:        emailLower, // email регистронезависимый
		Login:        loginLower, // логин регистронезависимый: храним в нижнем регистре
		PasswordHash: passwordHashStr,
	}

	err = uc.repo.Save(ctx, user)
	if err != nil {
		return dto.RegisterOutput{}, fmt.Errorf("save user: %w", err)
	}

	issued, err := uc.session.IssueSession(ctx, user.ID)
	if err != nil {
		return dto.RegisterOutput{}, fmt.Errorf("issue session: %w", err)
	}

	return dto.RegisterOutput{UserID: id, Session: issued}, nil
}
