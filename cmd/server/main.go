package main

import (
	"log"
	"net/http"

	"github.com/tahadeh2010/realtime-terminal-collab/internal/transport"
)

func main() {
	http.HandleFunc("/ws", transport.HandleWebSocket)

	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
