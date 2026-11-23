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
	repo      outbound.LoginUserRepository
	hasher    outbound.PasswordHasher
	dummyHash string
	session   *SessionIssuer
}

func NewLoginUseCase(
	repo outbound.LoginUserRepository,
	hasher outbound.PasswordHasher,
	sessions *SessionIssuer,
) (*LoginUseCase, error) {
	dummyHash, err := hasher.Hash("dummy-password")
	if err != nil {
		return nil, fmt.Errorf("new login use case: hash dummy password: %w", err)
	}

	return &LoginUseCase{
		repo:      repo,
		hasher:    hasher,
		dummyHash: dummyHash,
		session:   sessions,
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
		user, err = uc.repo.GetByEmail(ctx, domain.NormalizeEmail(input.Identifier))
	} else {
		user, err = uc.repo.GetByLogin(ctx, domain.NormalizeLogin(input.Identifier))
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

	issued, err := uc.session.IssueSession(ctx, user.ID)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("issue session: %w", err)
	}

	return dto.LoginOutput{
		AccessToken:      issued.AccessToken,
		RefreshToken:     issued.RefreshToken,
		AccessExpiresIn:  issued.AccessExpiresIn,
		RefreshExpiresIn: issued.RefreshExpiresIn,
	}, nil
}

type SessionIssuer struct {
	tokenIssuer outbound.TokenIssuer
	sessionRepo outbound.SessionRepository
	refreshTTL  time.Duration
}

func NewSessionIssuer(
	tokenIssuer outbound.TokenIssuer,
	sessionRepo outbound.SessionRepository,
	refreshTTL time.Duration,
) *SessionIssuer {

	return &SessionIssuer{
		tokenIssuer: tokenIssuer,
		sessionRepo: sessionRepo,
		refreshTTL:  refreshTTL,
	}

}
