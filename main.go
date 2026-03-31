package main

import (
	"log"
	"net/http"
	"time"

	"context"

	"github.com/W1zard70r/Ancord/internal/config"
	"github.com/W1zard70r/Ancord/internal/delivery/api"
	"github.com/W1zard70r/Ancord/internal/delivery/ws"
	"github.com/W1zard70r/Ancord/internal/repopsitory/postgres"
	"github.com/W1zard70r/Ancord/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("Не удалось создать пул подключений: %v", err)
	}

	userRepo := postgres.NewUserRepository(pool)
	chatRepo := postgres.NewChatRepository(pool)
	msgRepo := postgres.NewMessageRepository(pool)

	userUC := usecase.NewUserUseCase(userRepo)
	chatUC := usecase.NewChatUseCase(chatRepo)
	msgUC := usecase.NewMessageUseCase(msgRepo, chatUC)

	h := api.NewHandler(userUC, chatUC, msgUC, cfg.JWTKey)
	r := chi.NewRouter()

	// запуск websocket hub
	hub := ws.NewHub(msgUC, chatUC)
	go hub.Run()
	wsHandler := ws.NewWSHandler(hub)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Твой фронт
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	// Публичные
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.Login)
	// r.Post("/unlogin", h.Login)

	// Приватные
	r.Group(func(r chi.Router) {
		r.Use(api.AuthMiddleware([]byte(cfg.JWTKey)))
		r.Post("/chats", h.CreateChat)
		r.Post("/chats/{chat_id}/members", h.AddMember)
		r.Post("/message", h.SendMessage)
		r.Get("/chat/{chat_id}", h.GetChatByID)
		r.Get("/messages/{chat_id}", h.GetMessagesByChatID)
		r.Get("/message/{message_id}", h.GetMessageByID)
		r.Get("/ws", wsHandler.ServeWS)
	})

	chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s", method, route)
		return nil
	})

	log.Println("Сервер запущен на :8080")

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatalf("Сервер упал с ошибкой: %v", err)
	}

}
