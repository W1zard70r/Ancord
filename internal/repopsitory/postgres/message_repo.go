package postgres

import (
	"context"

	"github.com/W1zard70r/Ancord/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) domain.MessageRepository {
	return &messageRepository{pool: pool}
}

func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	query := `INSERT INTO messages id, chat_id, user_id, content, created_at VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, message.ID, message.ChatID, message.UserID, message.Content, message.CreatedAt)
	return err
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	query := `SELECT id, chat_id, user_id, content, created_at FROM messages WHERE id = $1`
	message := &domain.Message{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&message.ID, &message.ChatID, &message.UserID, &message.Content, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (r *messageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit int) ([]*domain.Message, error) {
	query := `SELECT id, chat_id, user_id, content, created_at 
				FROM messages 
				WHERE chat_id = $1 
				ORDER BY created_at DESC 
				LIMIT $2`
	rows, err := r.pool.Query(ctx, query, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		message := &domain.Message{}
		err := rows.Scan(&message.ID, &message.ChatID, &message.UserID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
