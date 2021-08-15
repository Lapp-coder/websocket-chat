package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var indexTemplate *template.Template

type handler struct {
	hub *hub
}

func newHandler(hub *hub) *handler {
	return &handler{hub: hub}
}

func (h *handler) initRoutes() {
	http.HandleFunc("/", h.index)
	http.HandleFunc("/chat", h.chat)
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	if err := indexTemplate.Execute(w, r.Host); err != nil {
		http.Error(w, "failed to load a home page", http.StatusInternalServerError)
		return
	}
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

	go conn.writeMessages()
	conn.readMessages()
}
