package main

import (
	"log"
	"net/http"

	"github.com/tahadeh2010/realtime-terminal-collab/internal/application"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/infrastructure"
	"github.com/tahadeh2010/realtime-terminal-collab/internal/transport"
)

func main() {
	store := infrastructure.NewMemoryStore()
	ptyManager := infrastructure.NewPTYManager()
	sm := application.NewSessionManager(store, ptyManager)
	cm := application.NewConnectionManager()

	ptyStreamer := transport.NewPTYStreamer(cm)
	go ptyStreamer.WatchSessions(sm)

	server := transport.NewServer(sm, cm)
	http.HandleFunc("/ws", server.HandleWebSocket)

	sessionHandler := transport.NewSessionHandler(sm)
	http.HandleFunc("/session", sessionHandler.CreateSession)

	log.Println("server starting on :5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
