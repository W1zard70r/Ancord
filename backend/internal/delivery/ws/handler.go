package ws

import (
	"log"
	"net/http"

	"github.com/W1zard70r/Ancord/internal/delivery/api"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //TODO: CSRF - защита
}

type WSHandler struct {
	hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// получение инфы с middleware
	val := r.Context().Value(api.UserIDKey)
	if val == nil {
		http.Error(w, "Unauthorized (no user ID in context)", http.StatusUnauthorized)
		return
	}
	userID := val.(uuid.UUID)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	// делаем клиента
	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
	}
	// регаем в хабе
	h.hub.register <- client

	// запускаем горутины на чтение и запись
	go client.writePump()
	go client.readPump()
}
