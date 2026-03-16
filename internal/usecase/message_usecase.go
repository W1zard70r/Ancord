package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/W1zard70r/Ancord/internal/domain"
	"github.com/google/uuid"
)

type MessageUseCase interface {
	Send(ctx context.Context, user_id uuid.UUID, chat_id uuid.UUID, content string) (*domain.Message, error)
}

type messageUseCase struct {
	messageRepo domain.MessageRepository
	chatUC      ChatUseCase
}

func NewMessageUseCase(repo domain.MessageRepository, chatUC ChatUseCase) MessageUseCase {
	return &messageUseCase{
		messageRepo: repo,
		chatUC:      chatUC,
	}
}

func (uc *messageUseCase) Send(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, content string) (*domain.Message, error) {
	// user, err := GetByID() Проверка на то, есть ли юзер-отправитель
	chat, err := uc.chatUC.GetChat(ctx, chatID) // проверка на чат
	if err != nil {
		return nil, err
	} else if chat == nil {
		return nil, fmt.Errorf("chat with id %s not found", chatID.String())
	}
	message := &domain.Message{
		ID:        uuid.Must(uuid.NewV7()),
		ChatID:    chatID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	err = uc.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
