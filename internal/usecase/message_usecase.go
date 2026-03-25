package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/W1zard70r/Ancord/internal/domain"
	"github.com/google/uuid"
)

type MessageUseCase interface {
	Send(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, content string) (*domain.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID, limit int) ([]*domain.Message, error)
	GetMessage(ctx context.Context, messageID uuid.UUID) (*domain.Message, error)
}

type messageUseCase struct {
	messageRepo domain.MessageRepository
	chatUC      ChatUseCase
	userUC      UserUseCase
}

func NewMessageUseCase(repo domain.MessageRepository, chatUC ChatUseCase) MessageUseCase {
	return &messageUseCase{
		messageRepo: repo,
		chatUC:      chatUC,
	}
}

func (uc *messageUseCase) Send(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, content string) (*domain.Message, error) {
	// user, err := uc.userUC.GetByID(ctx, userID) //Проверка на то, есть ли юзер-отправитель
	// if err != nil {
	// 	return nil, err
	// } else if user == nil {
	// 	return nil, fmt.Errorf("chat with id %s not found", chatID.String())
	// } // пока этого интерфейса нет
	// TODO: Сделать проверку, что пользователь в чате
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

func (uc *messageUseCase) GetMessages(ctx context.Context, chatID uuid.UUID, limit int) ([]*domain.Message, error) {
	return uc.messageRepo.GetByChatID(ctx, chatID, limit)
}

func (uc *messageUseCase) GetMessage(ctx context.Context, messageID uuid.UUID) (*domain.Message, error) {
	return uc.messageRepo.GetByID(ctx, messageID)
}
