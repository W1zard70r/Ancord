package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/W1zard70r/Ancord/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	userUC    usecase.UserUseCase
	chatUC    usecase.ChatUseCase
	messageUC usecase.MessageUseCase
}

func NewHandler(u usecase.UserUseCase, c usecase.ChatUseCase, m usecase.MessageUseCase) *Handler {
	return &Handler{
		userUC:    u,
		chatUC:    c,
		messageUC: m,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	// POST (Запись)
	r.Post("/register", h.RegisterUser)
	r.Post("/chats", h.CreateChat)
	r.Post("/message", h.SendMessage)

	// GET (Чтение)
	r.Get("/chat/{chat_id}", h.GetChatByID)
	r.Get("/message/{message_id}", h.GetMessageByID)
	r.Get("/messages/{chat_id}", h.GetMessagesByChatID)

	return r
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userUC.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat, err := h.chatUC.CreateChat(r.Context(), req.Name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ChatID  uuid.UUID `json:"chat_id"`
		UserID  uuid.UUID `json:"user_id"`
		Content string    `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.messageUC.Send(r.Context(), req.UserID, req.ChatID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetChatByID(w http.ResponseWriter, r *http.Request) {
	chatID, err := uuid.Parse(chi.URLParam(r, "chat_id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	chat, err := h.chatUC.GetChat(r.Context(), chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if chat == nil {
		http.Error(w, "Chat not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

func (h *Handler) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	messageID, err := uuid.Parse(chi.URLParam(r, "message_id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	message, err := h.messageUC.GetMessage(r.Context(), messageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if message == nil {
		http.Error(w, "Chat not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *Handler) GetMessagesByChatID(w http.ResponseWriter, r *http.Request) {
	chatID, err := uuid.Parse(chi.URLParam(r, "chat_id"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)

	if limit <= 0 {
		limit = 50
	}
	chat, err := h.chatUC.GetChat(r.Context(), chatID)
	if chat == nil {
		http.Error(w, "Chat not Found", http.StatusNotFound)
		return
	}

	messages, err := h.messageUC.GetMessages(r.Context(), chatID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
