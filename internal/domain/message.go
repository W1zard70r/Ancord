package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	ChatID    uuid.UUID `json:"chat_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   uuid.UUID `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageRepository interface {
	Create(message *Message) error
	GetByID(id uuid.UUID) (*Message, error)
	GetByChatID(chatID uuid.UUID, limit int) ([]*Message, error)
}
