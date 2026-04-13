package usecase

import (
	"context"
	"fmt"
	"time"

	"ancord-voice/backend/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.User, error)
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
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("Пользователь с таким username уже сущетсвует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Ошибка хэширования пароля: %w", err)
	}

	user = &domain.User{
		ID:           uuid.Must(uuid.NewV7()),
		Username:     username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *userUseCase) Login(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("Пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("неправильный пароль")
	}
	return user, nil
}
