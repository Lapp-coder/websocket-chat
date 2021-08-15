package main

// hub является центральным узлом (чатом), который обрабатывает все соединения
type hub struct {
	connections map[*connection]struct{}
	broadcast   chan []byte
	register    chan *connection
	unregister  chan *connection
}

func newHub() *hub {
	return &hub{
		connections: make(map[*connection]struct{}),
		broadcast:   make(chan []byte, 10000),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

// listen прослушивает новые соединения, закрытые соединения и сообщения из канала broadcast
func (h *hub) listen() {
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
				select {
				case conn.send <- message:
				default:
					delete(h.connections, conn)
					close(conn.send)
				}
			}
		}
	}
}
