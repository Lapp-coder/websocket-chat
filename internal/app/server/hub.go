package server

// Hub является центральным узлом (чатом), который обрабатывает все соединения.
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

// Connection возращает ссылку на соединения и булево значение о том,
// существует ли такое соединение в Hub.
func (h Hub) Connection(id string) (conn *connection, exists bool) {
	conn, exists = h.connections[id]
	return
}

// Listen прослушивает сигнал на создание нового соеднинения, закрытия соединения,
// а также сообщения из канала broadcast, куда попадают все публичные сообщения.
func (h *Hub) Listen() {
	for {
		select {
		case conn := <-h.register:
			h.registerConn(conn)
		case conn := <-h.unregister:
			h.unregisterConn(conn)
		case message := <-h.broadcast:
			h.sendMessageAllConnections(message)
		}
	}
}

func (h *Hub) registerConn(conn *connection) {
	h.connections[conn.id] = conn
}

func (h *Hub) unregisterConn(conn *connection) {
	if _, ok := h.connections[conn.id]; ok {
		delete(h.connections, conn.id)
		close(conn.send)
	}
}

func (h *Hub) sendMessageAllConnections(message []byte) {
	for _, conn := range h.connections {
		select {
		case conn.send <- message:
		default:
			delete(h.connections, conn.id)
			close(conn.send)
		}
	}
}
