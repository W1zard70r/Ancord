package postgres

import (
	"context"

	"github.com/W1zard70r/Ancord/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, username, password_hash, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Username, user.PasswordHash, user.CreatedAt)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil // Возвращаем nil, nil - это значит "не найдено"
	}
	if err != nil {
		return nil, err // Это реальная ошибка БД
	}
	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil // Возвращаем nil, nil - это значит "не найдено"
	}
	if err != nil {
		return nil, err // Это реальная ошибка БД
	}
	return user, nil
}
