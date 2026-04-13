package usecase

import (
	"context"
	"errors"
	"time"

	"ancord-voice/backend/internal/domain"

	"github.com/google/uuid"
)

type ChatUseCase interface {
	CreateChat(ctx context.Context, userID uuid.UUID, name string) (*domain.Chat, error)
	GetChat(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
	GetChatByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error)
	IsUserMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
	AddUserToChat(ctx context.Context, chatID, inviterID, targetUserID uuid.UUID) error
}

type chatUseCase struct {
	chatRepo domain.ChatRepository
}

func NewChatUseCase(repo domain.ChatRepository) ChatUseCase {
	return &chatUseCase{
		chatRepo: repo,
	}
}

func (uc *chatUseCase) CreateChat(ctx context.Context, userID uuid.UUID, name string) (*domain.Chat, error) {
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

	err = uc.chatRepo.AddMember(ctx, chat.ID, userID)
	return chat, nil
}

func (uc *chatUseCase) GetChat(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	return uc.chatRepo.GetByID(ctx, id)
}

func (uc *chatUseCase) GetChatByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error) {
	return uc.chatRepo.GetByUserID(ctx, userID)
}

func (uc *chatUseCase) IsUserMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	return uc.chatRepo.IsMember(ctx, chatID, userID)
}

func (uc *chatUseCase) AddUserToChat(ctx context.Context, chatID, inviterID, targetUserID uuid.UUID) error {
	isInviter, err := uc.chatRepo.IsMember(ctx, chatID, inviterID)
	if err != nil {
		return err
	}
	if !isInviter {
		return errors.New("user have no ability to invite")
	}

	isTarget, err := uc.chatRepo.IsMember(ctx, chatID, targetUserID)
	if err != nil {
		return err
	}
	if isTarget {
		return errors.New("user already in chat")
	}

	return uc.chatRepo.AddMember(ctx, chatID, targetUserID)
}
