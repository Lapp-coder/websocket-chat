package main

type hub struct {
	connections map[*connection]struct{}
	broadcast   chan []byte
	register    chan *connection
	unregister  chan *connection
}

func newHub() hub {
	return hub{
		connections: make(map[*connection]struct{}),
		broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

func (h *hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn] = struct{}{}
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				close(conn.send)
			}
		case message := <-h.broadcast:
			for conn := range h.connections {
				conn.send <- message
			}
		}
	}
}
