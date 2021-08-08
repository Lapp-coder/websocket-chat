package main

import "github.com/gorilla/websocket"

type connection struct {
	ws   *websocket.Conn
	send chan []byte
	hub  *hub
}

func newConnection(ws *websocket.Conn, hub *hub) *connection {
	return &connection{
		ws:   ws,
		send: make(chan []byte, 256),
		hub:  hub,
	}
}

func (c *connection) readMessages() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		c.hub.broadcast <- message
	}
	c.ws.Close()
}

func (c *connection) writeMessages() {
	for message := range c.send {
		if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
	c.ws.Close()
}
