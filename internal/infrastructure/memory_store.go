package infrastructure

import (
	"fmt"

	"github.com/tahadeh2010/realtime-terminal-collab/internal/application"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/domain"
)

type MemoryStore struct {
	sessions map[string]*domain.Session
}

var _ application.SessionStore = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions: make(map[string]*domain.Session),
	}
}

func (m *MemoryStore) Create(session *domain.Session) error {
	m.sessions[session.ID] = session
	return nil
}

func (m *MemoryStore) Get(id string) (*domain.Session, error) {
	session, ok := m.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session %s not found", id)
	}
	return session, nil
}

func (m *MemoryStore) Delete(id string) error {
	if _, ok := m.sessions[id]; !ok {
		return fmt.Errorf("session %s not found", id)
	}
	delete(m.sessions, id)
	return nil
}
