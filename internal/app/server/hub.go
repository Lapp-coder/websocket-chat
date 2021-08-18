package server

// Hub является центральным узлом (чатом), который обрабатывает все соединения
type Hub struct {
	connections map[string]*connection
	broadcast   chan []byte
	register    chan *connection
	unregister  chan *connection
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]*connection),
		broadcast:   make(chan []byte, 10000),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

// Listen прослушивает новые соединения, закрытые соединения и сообщения из канала broadcast
func (h *Hub) Listen() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn.id] = conn
		case conn := <-h.unregister:
			if _, ok := h.connections[conn.id]; ok {
				delete(h.connections, conn.id)
				close(conn.send)
			}
		case message := <-h.broadcast:
			for _, conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					delete(h.connections, conn.id)
					close(conn.send)
				}
			}
		}
	}
}
