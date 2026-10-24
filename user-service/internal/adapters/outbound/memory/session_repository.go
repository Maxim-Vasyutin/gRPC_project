// ВРЕМЕННЫЙ адаптер: хранит сессии в памяти процесса.
// Нужен, чтобы Login можно было собрать в app, пока нет Postgres-версии.
// После рестарта сервера все сессии пропадают.
// Когда появится PostgresSessionRepository — в app меняется одна строка.

package memory

import (
	"context"
	"sync"

	"user-service/internal/core/corerrors"
	"user-service/internal/core/domain"
	"user-service/internal/core/ports/outbound"
)

var _ outbound.SessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	mu       sync.Mutex                // защищает map: gRPC обрабатывает запросы в разных горутинах
	sessions map[string]domain.Session // ключ — refresh-токен
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{sessions: make(map[string]domain.Session)}
}

func (r *SessionRepository) Store(_ context.Context, s domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[s.RefreshToken] = s
	return nil
}

func (r *SessionRepository) Replace(_ context.Context, oldRefreshToken string, newSession domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[oldRefreshToken]; !ok {
		return corerrors.ErrSessionNotFound
	}
	delete(r.sessions, oldRefreshToken)
	r.sessions[newSession.RefreshToken] = newSession
	return nil
}

func (r *SessionRepository) FindByRefreshToken(_ context.Context, refreshToken string) (domain.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.sessions[refreshToken]
	if !ok {
		return domain.Session{}, corerrors.ErrSessionNotFound
	}
	return s, nil
}

func (r *SessionRepository) RevokeByRefreshToken(_ context.Context, refreshToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[refreshToken]; !ok {
		return corerrors.ErrSessionNotFound
	}
	delete(r.sessions, refreshToken)
	return nil
}
