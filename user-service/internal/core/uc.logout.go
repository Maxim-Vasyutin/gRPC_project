package core

import (
	"context"
	"fmt"
)

type LogoutInput struct {
	RefreshToken string
}

type LogoutUseCase struct {
	sessionRepo SessionRepository
}

func NewLogoutUseCase(
	sessionRepo SessionRepository,
) LogoutUseCase {
	return LogoutUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutUseCase) Logout(ctx context.Context, input LogoutInput) error {
	//Можно отдельно найти сессию, если нужно.
	//Метод Revoke и так её ищет дальше
	/*
		_, err := uc.sessionRepo.FindByRefreshToken(ctx, input.RefreshToken)
		if err != nil {
			return fmt.Errorf("find session: %w", err)
		}
	*/

	err := uc.sessionRepo.RevokeByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
