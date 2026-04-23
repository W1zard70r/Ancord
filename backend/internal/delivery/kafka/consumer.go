package kafka

import (
	"backend/internal/delivery/ws"
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
)

type VoiceEvent struct {
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id"`
	Action string `json:"action"`
}

func StartVoiceEventConsumer(hub *ws.Hub) {
	broker := os.Getenv("KAFKA_BROKER") // переделать в config.KAFKA_BROKER

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "voice-events",
		GroupID:  "backend-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	slog.Info("Kafka consumer started", slog.String("topic", "voice-events"))

	go func() {
		defer reader.Close()
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				slog.Error("Kafka read Error", slog.String("error", err.Error()))
				continue
			}
			var event VoiceEvent
			err = json.Unmarshal(m.Value, &event)
			if err != nil {
				slog.Error("message value Error", slog.String("error", err.Error()))
				continue
			}
			slog.Info(
				"Recieved voice event from Kafka",
				slog.String("user_id", event.UserID),
				slog.String("chat_id", event.ChatID),
			)

			chatID, err := uuid.Parse(event.ChatID)
			if err != nil {
				slog.Error("chat_id validate Error", slog.String("chat_id", chatID.String()), slog.String("error", err.Error()))
				continue
			}
			wsMessage := map[string]interface{}{
				"type": "voice_status",
				"data": event,
			}

			msgBytes, _ := json.Marshal(wsMessage)
			hub.BroadcastToChat(chatID, msgBytes)
		}
	}()
}
