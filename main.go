package main

import (
	"log"
	"net/http"
	"time"

	"context"

	"github.com/W1zard70r/Ancord/internal/config"
	"github.com/W1zard70r/Ancord/internal/delivery/api"
	"github.com/W1zard70r/Ancord/internal/repopsitory/postgres"
	"github.com/W1zard70r/Ancord/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, _ := pgxpool.New(ctx, cfg.DBURL)

	userRepo := postgres.NewUserRepository(pool)
	chatRepo := postgres.NewChatRepository(pool)
	msgRepo := postgres.NewMessageRepository(pool)

	userUC := usecase.NewUserUseCase(userRepo)
	chatUC := usecase.NewChatUseCase(chatRepo)
	msgUC := usecase.NewMessageUseCase(msgRepo, chatUC)

	h := api.NewHandler(userUC, chatUC, msgUC)
	r := chi.NewRouter()

	r.Mount("/api/v1", h.Routes())

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":"+cfg.Port, r)
}
