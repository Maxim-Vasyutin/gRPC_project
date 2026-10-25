package application

import (
	"context"
	"errors"
	"fmt"

	"user-service/internal/core/corerrors"
	"user-service/internal/core/dto"
	"user-service/internal/core/ports/outbound"
)

type LogoutUseCase struct {
	sessionRepo outbound.SessionRepository
}

func NewLogoutUseCase(
	sessionRepo outbound.SessionRepository,
) LogoutUseCase {
	return LogoutUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutUseCase) Logout(ctx context.Context, input dto.LogoutInput) error {
	//Можно отдельно найти сессию, если нужно.
	//Метод Revoke и так её ищет дальше
	/*
		_, err := uc.sessionRepo.FindByRefreshToken(ctx, input.RefreshToken)
		if err != nil {
			return fmt.Errorf("find session: %w", err)
		}
	*/

	err := uc.sessionRepo.RevokeByRefreshToken(ctx, input.RefreshToken)
	if errors.Is(err, corerrors.ErrSessionNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
