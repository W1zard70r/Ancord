package postgres

import (
	"context"

	"ancord-voice/backend/internal/domain"

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

func (r *chatRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error) {
	query := `SELECT chat_id, name, created_at 
				FROM chat_members 
				JOIN chats 
				ON chat_members.chat_id = chats.id
				WHERE user_id = $1 
				ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]*domain.Chat, 0)
	for rows.Next() {
		chat := &domain.Chat{}
		err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *chatRepository) IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, chatID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *chatRepository) AddMember(ctx context.Context, chatID, userID uuid.UUID) error {
	query := `INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, chatID, userID)
	return err
}
