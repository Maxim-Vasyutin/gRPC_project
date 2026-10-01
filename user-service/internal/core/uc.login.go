package core

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type LoginUseCase struct {
	repo        LoginUserRepository
	hasher      PasswordService
	tokenIssuer TokenIssuer
	sessionRepo SessionRepository
}

func NewLoginUseCase(
	repo LoginUserRepository,
	hasher PasswordService,
	tokenIssuer TokenIssuer,
	sessionRepo SessionRepository,
) LoginUseCase {
	return LoginUseCase{
		repo:        repo,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
		sessionRepo: sessionRepo,
	}
}

func (uc *LoginUseCase) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	//Вытягиваю пользователя
	user, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("get user by email: %w", err)
	}

	//Верификация введёного пароля и сохраннного хэша
	err = uc.hasher.Verify(input.Password, user.PasswordHash)
	if err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}

	//Присваивание токена сессии и рефреш-токена
	tokens, err := uc.tokenIssuer.Issue(user.ID)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("issue tokens: %w", err)
	}

	//создание сессии
	session := Session{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
	}

	//Сохраняем сессию
	err = uc.sessionRepo.Store(ctx, session)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("store session: %w", err)
	}

	return LoginOutput{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken}, nil
}
