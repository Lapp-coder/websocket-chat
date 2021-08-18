package server

import (
	"errors"
	"fmt"
	"github.com/Lapp-coder/websocket-chat/pkg/jrpc"
	"github.com/gorilla/websocket"
	"net/http"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) InitRoutes() {
	http.HandleFunc("/chat", h.chat)
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	uniqueId int
)

func (h *Handler) chat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	uniqueId++
	id := strconv.Itoa(uniqueId)
	if err = ws.WriteMessage(websocket.TextMessage, []byte(id)); err != nil {
		return
	}

	conn := newConnection(id, ws, h.hub)
	conn.hub.register <- conn
	defer func() {
		conn.hub.unregister <- conn
	}()

	go conn.writeMessages()

	jsonrpc.ServeConn(ws.UnderlyingConn())
}

func (h *Handler) SendMessage(args *jrpc.SendMessageArgs, result *string) error {
	switch args.Ids {
	case "*":
		h.hub.broadcast <- []byte(fmt.Sprintf("%s: %s", args.ID, args.Message))
	case "echo":
		conn, ok := h.hub.connections[args.ID]
		if ok {
			conn.send <- []byte(args.Message)
		}
	default:
		ids := strings.Split(args.Ids, ", ")
		if len(ids) == 0 {
			return errors.New("ids is empty")
		}

		for _, id := range ids {
			conn, ok := h.hub.connections[id]
			if ok {
				conn.send <- []byte(fmt.Sprintf("%s: %s", args.ID, args.Message))
			}
		}
	}

	*result = "OK"
	return nil
}

func (h *Handler) GetMessages(args *jrpc.GetMessagesArgs, result *[]string) error {
	conn, ok := h.hub.connections[args.ID]
	if ok {
		for key, message := range conn.unreadMessages {
			*result = append(*result, message)
			delete(conn.unreadMessages, key)
		}
	}

	return nil
}
