package usecase

import (
	"context"
	"fmt"
	"time"

	"ancord-voice/backend/internal/domain"

	"github.com/google/uuid"
)

type MessageUseCase interface {
	Send(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, content string) (*domain.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID, limit int) ([]*domain.Message, error)
	GetMessage(ctx context.Context, messageID uuid.UUID) (*domain.Message, error)
}

type messageUseCase struct {
	messageRepo domain.MessageRepository
	chatRepo    domain.ChatRepository
	chatUC      ChatUseCase
}

func NewMessageUseCase(messRepo domain.MessageRepository, chatUC ChatUseCase) MessageUseCase {
	return &messageUseCase{
		messageRepo: messRepo,
		chatUC:      chatUC,
	}
}

func (uc *messageUseCase) Send(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, content string) (*domain.Message, error) {

	//проверка на то, что пользователь в чате
	isMember, err := uc.chatUC.IsUserMember(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("вы не состоите в этом чате")
	}
	//создаём соообщение
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

func (uc *messageUseCase) GetMessages(ctx context.Context, chatID uuid.UUID, limit int) ([]*domain.Message, error) {
	return uc.messageRepo.GetByChatID(ctx, chatID, limit)
}

func (uc *messageUseCase) GetMessage(ctx context.Context, messageID uuid.UUID) (*domain.Message, error) {
	return uc.messageRepo.GetByID(ctx, messageID)
}
