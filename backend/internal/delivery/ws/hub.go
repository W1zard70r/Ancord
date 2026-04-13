package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"backend/internal/usecase"

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
	Type    string      `json:"type"`
	ChatID  uuid.UUID   `json:"chat_id"`
	Content string      `json:"content,omitempty"` // Для входящего текста от клиента
	Message interface{} `json:"message,omitempty"` // Для исходящего объекта (из БД)
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
		slog.Debug("Closing readPump", slog.String("user_id", c.UserID.String()))
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

			// получаем и кидаем историю
			messages, err := c.Hub.msgUC.GetMessages(context.Background(), action.ChatID, 50)
			if err == nil && len(messages) > 0 {
				response := map[string]interface{}{
					"type": "history",
					"data": messages,
				}
				respJSON, _ := json.Marshal(response)
				c.Send <- respJSON

			}

		case "message":
			// проверяем подписку
			c.Hub.mu.RLock() // Используй RLock для чтения (это быстрее и безопаснее)
			isSubscribed := c.Hub.subscriptions[action.ChatID][c.UserID]
			c.Hub.mu.RUnlock() // Сразу разблокировали!

			if !isSubscribed {
				c.Send <- []byte(`{"type":"error", "message":"Join chat before sending messages"}`)
				continue
			}

			// пробудем записать сообщение
			msg, err := c.Hub.msgUC.Send(context.Background(), c.UserID, action.ChatID, action.Content)
			if err != nil {
				c.Send <- []byte(`{"type":"error", "message":"` + err.Error() + `"}`)
				continue // Просто игнорируем это сообщение
			}

			// кидаем онлайн пользователям-подписчикам
			c.Hub.broadcast <- &WSAction{
				Type:    "message",
				ChatID:  action.ChatID,
				Message: msg,
			}
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
			slog.Info("Client connected", slog.String("user_id", client.UserID.String()))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)

				for chatID, subs := range h.subscriptions {
					if _, isSubbed := subs[client.UserID]; isSubbed {
						delete(subs, client.UserID)
						slog.Info("Client disconnected: user_id; chat_id", slog.String("user_id", client.UserID.String()), slog.String("chat_id", chatID.String()))
					}
				}
			}
			h.mu.Unlock()
		case action := <-h.broadcast:
			// 1. Формируем соощбение для отправки
			msgJSON, err := json.Marshal(action)
			if err != nil {
				slog.Error("Failed to marshal broadcast message", slog.String("error", err.Error()))
				continue
			}
			h.mu.RLock()
			// 2. Берем всех юзеров, кто в этом чате
			usersInChat := h.subscriptions[action.ChatID]
			slog.Debug("Broadcast message",
				slog.String("chat_id", action.ChatID.String()),
				slog.Int("target_users_count", len(usersInChat)),
			)
			// 3. Шлем только тем, кто сейчас онлайн И в этом чате
			for userID := range usersInChat {
				if client, ok := h.clients[userID]; ok {
					select {
					case client.Send <- msgJSON:
						// Успешно ушло
					default:
						// не получилось доставить всё
						slog.Warn("Client buffer full, dropping message", slog.String("user_id", userID.String()))
						// Можно принудительно отключить такого клиента
						// h.unregister <- client
					}
				} else {
					slog.Debug("User in subscriptions but not online", slog.String("user_id", userID.String()))
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
