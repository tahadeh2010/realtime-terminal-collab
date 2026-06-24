package transport

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/application"
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

	_, err = s.sm.GetSession(sessionID)
	if err != nil {
		log.Printf("session not found: %s", sessionID)
		conn.Close()
		return
	}

	clientID := uuid.New().String()
	s.cm.Register(clientID, conn)
	log.Printf("client %s joined session %s", clientID, sessionID)
}
