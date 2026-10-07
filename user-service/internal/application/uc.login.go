package application

import (
	"context"
	"fmt"

	"gRPC_project/user-service/internal/application/dto"
	"gRPC_project/user-service/internal/application/ports/outbound"
	"gRPC_project/user-service/internal/core"
)

type LoginUseCase struct {
	repo        outbound.LoginUserRepository
	hasher      outbound.PasswordHasher
	tokenIssuer outbound.TokenIssuer
	sessionRepo outbound.SessionRepository
}

func NewLoginUseCase(
	repo outbound.LoginUserRepository,
	hasher outbound.PasswordHasher,
	tokenIssuer outbound.TokenIssuer,
	sessionRepo outbound.SessionRepository,
) LoginUseCase {
	return LoginUseCase{
		repo:        repo,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
		sessionRepo: sessionRepo,
	}
}

func (uc *LoginUseCase) Login(ctx context.Context, input dto.LoginInput) (dto.LoginOutput, error) {
	//Вытягиваю пользователя
	user, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("get user by email: %w", err)
	}

	//Верификация введёного пароля и сохраннного хэша
	err = uc.hasher.Verify(input.Password, user.PasswordHash)
	if err != nil {
		return dto.LoginOutput{}, ErrInvalidCredentials
	}

	//Присваивание токена сессии и рефреш-токена
	tokens, err := uc.tokenIssuer.Issue(user.ID)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("issue tokens: %w", err)
	}

	//создание сессии
	session := core.Session{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
	}

	//Сохраняем сессию
	err = uc.sessionRepo.Store(ctx, session)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("store session: %w", err)
	}

	return dto.LoginOutput{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken}, nil
}
