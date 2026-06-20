package application

import (
	"github.com/google/uuid"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/domain"
)

type SessionManager struct {
	store SessionStore
}

func NewSessionManager(store SessionStore) *SessionManager {
	return &SessionManager{store: store}
}

func (s *SessionManager) CreateSession(hostID string) (*domain.Session, error) {
	session := &domain.Session{
		ID: uuid.New().String(),
		Host: &domain.Client{
			ID:   hostID,
			Role: domain.RoleHost,
		},
	}

	if err := s.store.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionManager) GetSession(id string) (*domain.Session, error) {
	return s.store.Get(id)
}

func (s *SessionManager) DeleteSession(id string) error {
	return s.store.Delete(id)
}
