package services

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type NotifyHub struct {
	mu      sync.RWMutex
	clients map[uint]map[*websocket.Conn]struct{}
}

func NewNotifyHub() *NotifyHub {
	return &NotifyHub{
		clients: make(map[uint]map[*websocket.Conn]struct{}),
	}
}

func (h *NotifyHub) Register(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[userID]; !ok {
		h.clients[userID] = make(map[*websocket.Conn]struct{})
	}
	h.clients[userID][conn] = struct{}{}
}

func (h *NotifyHub) Unregister(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[userID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.clients, userID)
		}
	}
	_ = conn.Close()
}

func (h *NotifyHub) SendToUser(userID uint, payload interface{}) {
	h.mu.RLock()
	connsMap := h.clients[userID]
	if len(connsMap) == 0 {
		h.mu.RUnlock()
		return
	}
	conns := make([]*websocket.Conn, 0, len(connsMap))
	for conn := range connsMap {
		conns = append(conns, conn)
	}
	h.mu.RUnlock()

	for _, conn := range conns {
		_ = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
		if err := conn.WriteJSON(payload); err != nil {
			h.Unregister(userID, conn)
		}
	}
}
