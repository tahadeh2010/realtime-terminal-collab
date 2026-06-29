package transport

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/application"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/domain"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	sm *application.SessionManager
	cm *application.ConnectionManager
}

func NewServer(sm *application.SessionManager, cm *application.ConnectionManager) *Server {
	return &Server{sm: sm, cm: cm}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionID")
	if sessionID == "" {
		http.Error(w, "sessionID required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	session, err := s.sm.GetSession(sessionID)
	if err != nil {
		log.Printf("session not found: %s", sessionID)
		conn.Close()
		return
	}

	clientID := uuid.New().String()
	s.cm.Register(clientID, conn)

	role := domain.RoleViewer
	if session.Host == nil {
		role = domain.RoleHost
		session.Host = &domain.Client{ID: clientID, Role: domain.RoleHost}
	}
	log.Printf("client %s joined session %s as %v", clientID, sessionID, role)

	s.readLoop(conn, sessionID, clientID, role)
}

func (s *Server) readLoop(conn *websocket.Conn, sessionID, clientID string, role domain.Role) {
	defer func() {
		s.cm.Unregister(clientID)
		conn.Close()
		log.Printf("client %s disconnected from session %s", clientID, sessionID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("read error: %v", err)
			}
			return
		}

		if role != domain.RoleHost {
			log.Printf("viewer %s attempted input, rejected", clientID)
			continue
		}

		pty, ok := s.sm.GetPTY(sessionID)
		if !ok {
			log.Printf("PTY not found for session %s", sessionID)
			return
		}

		if err := pty.Write(message); err != nil {
			log.Printf("PTY write error: %v", err)
			return
		}
	}
}
