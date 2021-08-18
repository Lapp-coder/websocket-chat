package server

import "github.com/gorilla/websocket"

type connection struct {
	id             string
	ws             *websocket.Conn
	send           chan []byte
	hub            *Hub
	unreadMessages map[int]string
}

func newConnection(id string, ws *websocket.Conn, hub *Hub) *connection {
	return &connection{
		id:             id,
		ws:             ws,
		send:           make(chan []byte, 256),
		hub:            hub,
		unreadMessages: make(map[int]string),
	}
}

func (c *connection) writeMessages() {
	var key int
	for message := range c.send {
		c.unreadMessages[key] = string(message)
		key++
	}
}
