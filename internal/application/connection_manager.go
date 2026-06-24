package application

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ID   string
	Conn *websocket.Conn
}

type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*Connection
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
	}
}

func (cm *ConnectionManager) Register(id string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.connections[id] = &Connection{
		ID:   id,
		Conn: conn,
	}
}

func (cm *ConnectionManager) Unregister(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.connections, id)
}

func (cm *ConnectionManager) GetConnection(id string) (*Connection, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conn, ok := cm.connections[id]
	return conn, ok
}

func (cm *ConnectionManager) Broadcast(message []byte) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for id, c := range cm.connections {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("broadcast error to %s: %v", id, err)
			c.Conn.Close()
		}
	}
}
