package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/W1zard70r/Ancord/internal/domain"
	"github.com/google/uuid"
)

type ChatUseCase interface {
	CreateChat(ctx context.Context, name string) (*domain.Chat, error)
	GetChat(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
}

type chatUseCase struct {
	chatRepo domain.ChatRepository
}

func NewChatUseCase(repo domain.ChatRepository) ChatUseCase {
	return &chatUseCase{
		chatRepo: repo,
	}
}

func (uc *chatUseCase) CreateChat(ctx context.Context, name string) (*domain.Chat, error) {
	if name == "" {
		return nil, errors.New("chat name cannot be empty")
	}

	chat := &domain.Chat{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      name,
		CreatedAt: time.Now(),
	}

	err := uc.chatRepo.Create(ctx, chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (uc *chatUseCase) GetChat(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	return uc.chatRepo.GetByID(ctx, id)
}
