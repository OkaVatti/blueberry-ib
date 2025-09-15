package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type WSHub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan WSMessage
	mutex     sync.Mutex
	upgrader  websocket.Upgrader
}

func NewWSHub() *WSHub {
	return &WSHub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan WSMessage),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *WSHub) Run() {
	for msg := range h.broadcast {
		h.mutex.Lock()
		for c := range h.clients {
			_ = c.WriteJSON(msg) // ignore error per-client for now
		}
		h.mutex.Unlock()
	}
}

func (h *WSHub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	// read pump (we don't expect messages, but keep connection alive)
	go func() {
		defer func() {
			h.mutex.Lock()
			delete(h.clients, conn)
			h.mutex.Unlock()
			conn.Close()
		}()
		for {
			if _, _, err := conn.NextReader(); err != nil {
				return
			}
		}
	}()
}
