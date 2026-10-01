package core

import (
	"context"
	"fmt"
)

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken  string
	RefreshToken string
}

type RefreshUseCase struct {
	sessionRepo SessionRepository
	tokenIssuer TokenIssuer
}

func NewRefreshUseCase(
	session SessionRepository,
	tokenIssuer TokenIssuer,
) RefreshUseCase {
	return RefreshUseCase{
		sessionRepo: session,
		tokenIssuer: tokenIssuer,
	}
}

func (uc *RefreshUseCase) Refresh(ctx context.Context, input RefreshInput) (RefreshOutput, error) {
	//Нахожу старую сессию
	oldSession, err := uc.sessionRepo.FindByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return RefreshOutput{}, fmt.Errorf("find refresh session: %w", err)
	}

	//Выпускаю новые токены
	tokens, err := uc.tokenIssuer.Issue(oldSession.UserID)
	if err != nil {
		return RefreshOutput{}, fmt.Errorf("issue tokens: %w", err)
	}

	//Собираю новую сессию
	newsession := Session{
		UserID:       oldSession.UserID,
		RefreshToken: tokens.RefreshToken,
	}

	//Заменяю старую сессию
	err = uc.sessionRepo.Replace(ctx, input.RefreshToken, newsession)
	if err != nil {
		return RefreshOutput{}, fmt.Errorf("replace refresh session: %w", err)
	}

	//Возвращаю новую пару токенов
	return RefreshOutput{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
