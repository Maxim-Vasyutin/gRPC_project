package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"user-service/internal/core/corerrors"
	"user-service/internal/core/domain"
	"user-service/internal/core/dto"
	"user-service/internal/core/ports/outbound"
)

// Хэш-заглушка (валидный SHA-256 формат в hex)
const MockHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

type LoginUseCase struct {
	repo        outbound.LoginUserRepository
	hasher      outbound.PasswordHasher
	tokenIssuer outbound.TokenIssuer
	sessionRepo outbound.SessionRepository
	dummyHash   string
	refreshTTL  time.Duration
}

func NewLoginUseCase(
	repo outbound.LoginUserRepository,
	hasher outbound.PasswordHasher,
	tokenIssuer outbound.TokenIssuer,
	sessionRepo outbound.SessionRepository,
	refreshTTL time.Duration,
) LoginUseCase {
	dummyHash, err := hasher.Hash("dummy-password")
	if err != nil {
		log.Printf("dummy hash not generated")
		return LoginUseCase{
			repo:        repo,
			hasher:      hasher,
			tokenIssuer: tokenIssuer,
			sessionRepo: sessionRepo,
			dummyHash:   MockHash,
			refreshTTL:  refreshTTL}
	}

	return LoginUseCase{
		repo:        repo,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
		sessionRepo: sessionRepo,
		dummyHash:   dummyHash,
		refreshTTL:  refreshTTL,
	}
}

func (uc *LoginUseCase) Login(ctx context.Context, input dto.LoginInput) (dto.LoginOutput, error) {
	//Вытягиваю пользователя
	user, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		//смотрим, что ошибка с ненайденным пользователем
		//и отправляем общую ошибку, чтобы было безопаснее
		if errors.Is(err, corerrors.ErrUserNotFound) {
			uc.hasher.Verify(input.Password, uc.dummyHash)
			return dto.LoginOutput{}, corerrors.ErrInvalidCredentials
		}
		return dto.LoginOutput{}, fmt.Errorf("get user by email: %w", err)
	}

	//Верификация введёного пароля и сохраннного хэша
	err = uc.hasher.Verify(input.Password, user.PasswordHash)
	if err != nil {
		return dto.LoginOutput{}, corerrors.ErrInvalidCredentials
	}

	//Присваивание токена сессии и рефреш-токена
	tokens, err := uc.tokenIssuer.Issue(user.ID)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("issue tokens: %w", err)
	}

	//создание сессии
	session := domain.Session{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(uc.refreshTTL),
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
