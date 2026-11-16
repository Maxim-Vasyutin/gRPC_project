package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"user-service/internal/core/corerrors"
	"user-service/internal/core/domain"
	"user-service/internal/core/dto"
	"user-service/internal/core/ports/outbound"
)

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
) (*LoginUseCase, error) {
	dummyHash, err := hasher.Hash("dummy-password")
	if err != nil {
		return nil, fmt.Errorf("new login use case: hash dummy password: %w", err)
	}

	return &LoginUseCase{
		repo:        repo,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
		sessionRepo: sessionRepo,
		dummyHash:   dummyHash,
		refreshTTL:  refreshTTL,
	}, nil
}

func (uc *LoginUseCase) Login(ctx context.Context, input dto.LoginInput) (dto.LoginOutput, error) {
	//Вытягиваю пользователя: есть "@" — ищем по email, иначе по логину.
	//Логин в базе хранится в нижнем регистре.
	var (
		user domain.User
		err  error
	)
	if strings.Contains(input.Identifier, "@") {
		user, err = uc.repo.GetByEmail(ctx, strings.ToLower(input.Identifier))
	} else {
		user, err = uc.repo.GetByLogin(ctx, strings.ToLower(input.Identifier))
	}
	if err != nil {
		//смотрим, что ошибка с ненайденным пользователем
		//и отправляем общую ошибку, чтобы было безопаснее
		if errors.Is(err, corerrors.ErrUserNotFound) {
			uc.hasher.Verify(input.Password, uc.dummyHash)
			return dto.LoginOutput{}, corerrors.ErrInvalidCredentials
		}
		return dto.LoginOutput{}, fmt.Errorf("get user by identifier: %w", err)
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
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresIn:  tokens.AccessTTL,
		RefreshExpiresIn: uc.refreshTTL,
	}, nil
}
