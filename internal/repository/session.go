package repository

import (
	"github.com/belyaevedu/remote-code-service/internal/domain"
)

func (r *Repository) CreateSession(session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clone := *session
	r.sessions[clone.SessionID] = &clone

	return nil
}

func (r *Repository) GetSession(sessionID string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, domain.ErrSessionNotFound
	}

	clone := *s
	return &clone, nil
}

func (r *Repository) DeleteSession(sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[sessionID]; !ok {
		return domain.ErrSessionNotFound
	}

	delete(r.sessions, sessionID)

	return nil
}
