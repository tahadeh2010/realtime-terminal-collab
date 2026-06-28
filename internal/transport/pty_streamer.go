package transport

import (
	"log"
	"sync"

	"github.com/tahadeh2010/realtime-terminal-collab/internal/application"
)

type PTYStreamer struct {
	cm      *application.ConnectionManager
	streams map[string]bool
	mu      sync.Mutex
}

func NewPTYStreamer(cm *application.ConnectionManager) *PTYStreamer {
	return &PTYStreamer{
		cm:      cm,
		streams: make(map[string]bool),
	}
}

func (s *PTYStreamer) WatchSessions(sm *application.SessionManager) {
	// Poll for new sessions with PTYs
	// In production, use channels/events
	ticker := make(chan struct{})
	go func() {
		for range ticker {
			s.checkSessions(sm)
		}
	}()

	// Initial check
	s.checkSessions(sm)
	// Block forever
	select {}
}

func (s *PTYStreamer) checkSessions(sm *application.SessionManager) {
	ids := sm.GetAllSessionIDs()
	for _, id := range ids {
		s.tryStream(id, sm)
	}
}

func (s *PTYStreamer) tryStream(sessionID string, sm *application.SessionManager) {
	s.mu.Lock()
	if s.streams[sessionID] {
		s.mu.Unlock()
		return
	}

	pty, ok := sm.GetPTY(sessionID)
	if !ok {
		s.mu.Unlock()
		return
	}

	s.streams[sessionID] = true
	s.mu.Unlock()

	go s.stream(sessionID, pty.Output())
	log.Printf("started streaming for session %s", sessionID)
}

func (s *PTYStreamer) stream(sessionID string, output <-chan []byte) {
	for data := range output {
		s.cm.Broadcast(data)
	}
	log.Printf("stream ended for session %s", sessionID)

	s.mu.Lock()
	delete(s.streams, sessionID)
	s.mu.Unlock()
}
