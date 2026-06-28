package application

import (
	"fmt"
	"testing"

	"github.com/tahadeh2010/realtime-terminal-collab/internal/domain"
)

type mockStore struct {
	sessions map[string]*domain.Session
}

func newMockStore() *mockStore {
	return &mockStore{sessions: make(map[string]*domain.Session)}
}

func (m *mockStore) Create(session *domain.Session) error {
	m.sessions[session.ID] = session
	return nil
}

func (m *mockStore) Get(id string) (*domain.Session, error) {
	s, ok := m.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session %s not found", id)
	}
	return s, nil
}

func (m *mockStore) Delete(id string) error {
	if _, ok := m.sessions[id]; !ok {
		return fmt.Errorf("session %s not found", id)
	}
	delete(m.sessions, id)
	return nil
}

type mockPTY struct {
	output chan []byte
}

func newMockPTY() *mockPTY {
	return &mockPTY{output: make(chan []byte, 1)}
}

func (m *mockPTY) Write(data []byte) error { return nil }
func (m *mockPTY) Output() <-chan []byte    { return m.output }
func (m *mockPTY) Close() error            { return nil }

type mockPTYProvider struct{}

func (m *mockPTYProvider) Spawn() (PTYInstance, error) {
	return newMockPTY(), nil
}

func (m *mockPTYProvider) Stop(inst PTYInstance) error {
	return inst.Close()
}

func TestCreateSession(t *testing.T) {
	manager := NewSessionManager(newMockStore(), &mockPTYProvider{})

	session, err := manager.CreateSession("host-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.ID == "" {
		t.Error("session ID should not be empty")
	}

	if session.Host == nil {
		t.Error("session host should not be nil")
	}

	if session.Host.ID != "host-1" {
		t.Errorf("expected host ID host-1, got %s", session.Host.ID)
	}
}

func TestGetSession(t *testing.T) {
	manager := NewSessionManager(newMockStore(), &mockPTYProvider{})

	created, err := manager.CreateSession("host-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fetched, err := manager.GetSession(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("expected session ID %s, got %s", created.ID, fetched.ID)
	}
}

func TestGetSessionNotFound(t *testing.T) {
	manager := NewSessionManager(newMockStore(), &mockPTYProvider{})

	_, err := manager.GetSession("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}
}

func TestDeleteSession(t *testing.T) {
	manager := NewSessionManager(newMockStore(), &mockPTYProvider{})

	created, err := manager.CreateSession("host-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = manager.DeleteSession(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = manager.GetSession(created.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestDeleteSessionNotFound(t *testing.T) {
	manager := NewSessionManager(newMockStore(), &mockPTYProvider{})

	err := manager.DeleteSession("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}
}
