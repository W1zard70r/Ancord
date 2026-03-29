package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/W1zard70r/Ancord/internal/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

type WSAction struct {
	Type    string    `json:"type"` // "join", "message"
	ChatID  uuid.UUID `json:"chat_id"`
	Content string    `json:"content"`
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var action WSAction
		json.Unmarshal(message, &action)

		switch action.Type {
		case "join":
			//проверка прав
			isMember, err := c.Hub.chatUC.IsUserMember(context.Background(), action.ChatID, c.UserID)
			if err != nil || !isMember {
				c.Send <- []byte(`{"type":"error", "message":"access denied"}`)
				continue
			}

			//подписываем
			c.Hub.mu.Lock()
			if c.Hub.subscriptions[action.ChatID] == nil {
				c.Hub.subscriptions[action.ChatID] = make(map[uuid.UUID]bool)
			}
			c.Hub.subscriptions[action.ChatID][c.UserID] = true
			c.Hub.mu.Unlock()

		case "message":
			// Сначала проверяем права!
			isMember, _ := c.Hub.chatUC.IsUserMember(context.Background(), action.ChatID, c.UserID)
			if !isMember {
				c.Send <- []byte(`{"type":"error", "message":"not a member"}`)
				continue
			}
			// Если ок — кидаем в хаб
			c.Hub.broadcast <- &action
		}
	}
}

type Hub struct {
	clients       map[uuid.UUID]*Client
	subscriptions map[uuid.UUID]map[uuid.UUID]bool

	broadcast  chan *WSAction
	register   chan *Client
	unregister chan *Client

	msgUC  usecase.MessageUseCase
	chatUC usecase.ChatUseCase

	mu sync.RWMutex
}

type MessagePayload struct {
	ChatID  uuid.UUID
	Message []byte
}

func NewHub(msgUC usecase.MessageUseCase, chatUC usecase.ChatUseCase) *Hub {
	return &Hub{
		clients:       make(map[uuid.UUID]*Client),
		subscriptions: make(map[uuid.UUID]map[uuid.UUID]bool),
		// Канал готов передавать указатели на WSAction
		broadcast:  make(chan *WSAction),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		msgUC:  msgUC,
		chatUC: chatUC,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, (client.UserID))
				close(client.Send)
			}
			h.mu.Unlock()
		case action := <-h.broadcast:

			h.mu.RLock()
			// 1. Берем всех юзеров, кто в этом чате
			usersInChat := h.subscriptions[action.ChatID]

			// 2. Формируем сообщение для отправки
			msgJSON, _ := json.Marshal(action)

			// 3. Шлем только тем, кто сейчас онлайн И в этом чате
			for userID := range usersInChat {
				if client, ok := h.clients[userID]; ok {
					client.Send <- msgJSON
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		client.Send <- message
	}
}
