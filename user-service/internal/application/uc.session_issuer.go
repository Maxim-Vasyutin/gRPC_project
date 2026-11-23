package application

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/core/domain"
	"user-service/internal/core/dto"
)

func (s *SessionIssuer) IssueSession(ctx context.Context, userID string) (dto.IssuedSession, error) {

	// ↓ выдача сессии
	//Присваивание токена сессии и рефреш-токена
	tokens, err := s.tokenIssuer.Issue(userID)
	if err != nil {
		return dto.IssuedSession{}, fmt.Errorf("issue tokens: %w", err)
	}

	//создание сессии
	session := domain.Session{
		UserID:       userID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTTL),
	}

	//Сохраняем сессию
	err = s.sessionRepo.Store(ctx, session)
	if err != nil {
		return dto.IssuedSession{}, fmt.Errorf("store session: %w", err)
	}

	return dto.IssuedSession{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresIn:  tokens.AccessTTL,
		RefreshExpiresIn: s.refreshTTL,
	}, nil

}
