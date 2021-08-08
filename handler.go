package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type handler struct {
	hub *hub
}

func newHandler(hub *hub) *handler {
	return &handler{hub: hub}
}

func (h *handler) initRoutes() {
	http.HandleFunc("/chat", h.chat)
}

func (h *handler) chat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := newConnection(ws, h.hub)
	conn.hub.register <- conn
	defer func() {
		conn.hub.unregister <- conn
	}()

	go conn.write()
	conn.read()
}
