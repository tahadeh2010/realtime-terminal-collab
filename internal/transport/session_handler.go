package transport

import (
	"encoding/json"
	"net/http"

	"github.com/tahadeh2010/realtime-terminal-collab/internal/application"
)

type SessionHandler struct {
	sm *application.SessionManager
}

func NewSessionHandler(sm *application.SessionManager) *SessionHandler {
	return &SessionHandler{sm: sm}
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := h.sm.CreateSession("host")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}
