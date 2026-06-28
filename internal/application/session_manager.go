package application

import (
	"sync"

	"github.com/google/uuid"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/domain"
)

type SessionManager struct {
	store      SessionStore
	ptyManager PTYProvider
	ptys       map[string]PTYInstance
	mu         sync.RWMutex
}

func NewSessionManager(store SessionStore, ptyManager PTYProvider) *SessionManager {
	return &SessionManager{
		store:      store,
		ptyManager: ptyManager,
		ptys:       make(map[string]PTYInstance),
	}
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

	pty, err := s.ptyManager.Spawn()
	if err != nil {
		s.store.Delete(session.ID)
		return nil, err
	}

	s.mu.Lock()
	s.ptys[session.ID] = pty
	s.mu.Unlock()

	return session, nil
}

func (s *SessionManager) GetSession(id string) (*domain.Session, error) {
	return s.store.Get(id)
}

func (s *SessionManager) DeleteSession(id string) error {
	s.mu.Lock()
	if pty, ok := s.ptys[id]; ok {
		s.ptyManager.Stop(pty)
		delete(s.ptys, id)
	}
	s.mu.Unlock()

	return s.store.Delete(id)
}

func (s *SessionManager) GetPTY(sessionID string) (PTYInstance, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pty, ok := s.ptys[sessionID]
	return pty, ok
}

func (s *SessionManager) GetAllSessionIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]string, 0, len(s.ptys))
	for id := range s.ptys {
		ids = append(ids, id)
	}
	return ids
}
