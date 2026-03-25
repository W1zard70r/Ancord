package postgres

import (
	"context"

	"github.com/W1zard70r/Ancord/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type chatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) domain.ChatRepository {
	return &chatRepository{pool: pool}
}

func (r *chatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	query := `INSERT INTO chats (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, chat.ID, chat.Name, chat.CreatedAt)
	return err
}

func (r *chatRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	query := `SELECT id, name, created_at FROM chats WHERE id = $1`
	chat := &domain.Chat{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil // Возвращаем nil, nil - это значит "не найдено"
	}
	if err != nil {
		return nil, err // Это реальная ошибка БД
	}
	return chat, nil
}
