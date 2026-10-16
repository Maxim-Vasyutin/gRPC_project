package application

import (
	"context"
	"fmt"

	"user-service/internal/core/domain"
	"user-service/internal/core/dto"
	"user-service/internal/core/ports/outbound"
)

type RefreshUseCase struct {
	sessionRepo outbound.SessionRepository
	tokenIssuer outbound.TokenIssuer
}

func NewRefreshUseCase(
	session outbound.SessionRepository,
	tokenIssuer outbound.TokenIssuer,
) RefreshUseCase {
	return RefreshUseCase{
		sessionRepo: session,
		tokenIssuer: tokenIssuer,
	}
}

func (uc *RefreshUseCase) Refresh(ctx context.Context, input dto.RefreshInput) (dto.RefreshOutput, error) {
	//Нахожу старую сессию
	oldSession, err := uc.sessionRepo.FindByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return dto.RefreshOutput{}, fmt.Errorf("find refresh session: %w", err)
	}

	//Выпускаю новые токены
	tokens, err := uc.tokenIssuer.Issue(oldSession.UserID)
	if err != nil {
		return dto.RefreshOutput{}, fmt.Errorf("issue tokens: %w", err)
	}

	//Собираю новую сессию
	newsession := domain.Session{
		UserID:       oldSession.UserID,
		RefreshToken: tokens.RefreshToken,
	}

	//Заменяю старую сессию
	err = uc.sessionRepo.Replace(ctx, input.RefreshToken, newsession)
	if err != nil {
		return dto.RefreshOutput{}, fmt.Errorf("replace refresh session: %w", err)
	}

	//Возвращаю новую пару токенов
	return dto.RefreshOutput{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
