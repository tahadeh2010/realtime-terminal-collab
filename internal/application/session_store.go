package application

import "github.com/tahadeh2010/realtime-terminal-collab/internal/domain"

type SessionStore interface {
	Create(session *domain.Session) error
	Get(id string) (*domain.Session, error)
	Delete(id string) error
}
