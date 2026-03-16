package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRepository interface {
	Create(ctx context.Context, chat *Chat) error
	GetByID(ctx context.Context, id uuid.UUID) (*Chat, error)
}
