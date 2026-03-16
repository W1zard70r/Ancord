package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/W1zard70r/Ancord/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserUseCase interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
}

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (uc *userUseCase) Register(ctx context.Context, username, password string) (*domain.User, error) {
	// логика проверки можно ли следать такого пользователя
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("Пользователь с таким username уже сущетсвует")
	}
	user = &domain.User{
		ID:           uuid.Must(uuid.NewV7()),
		Username:     username,
		PasswordHash: password,
		CreatedAt:    time.Now(),
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
