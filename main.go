package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/W1zard70r/Ancord/internal/repopsitory/postgres"
	"github.com/W1zard70r/Ancord/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := "postgres://postgres:secretpassword@127.0.0.1:5433/ancord_db?sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("БД не отвечает %v", err)
	}
	fmt.Println("Успешно подключились к БД")

	chatRepo := postgres.NewChatRepository(pool)
	chatUC := usecase.NewChatUseCase(chatRepo)
	fmt.Printf("Chat usecase initialized: %v\n", chatUC)

	ctxCreate, cancelCreate := context.WithTimeout(context.Background(), 2*time.Second)

}
