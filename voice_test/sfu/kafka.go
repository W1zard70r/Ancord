package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

type VoiceEvent struct {
	UserId string `json:"user_id"`
	ChatId string `json:"chat_id"`
	Action string `json:"action"`
}

func InitKafka() {
	broker := os.Getenv("KAFKA_BROKER")

	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "voice-events",
		Balancer: &kafka.LeastBytes{},
	}
	log.Printf("Kafka producer initialized for broker: %s\n", broker)
}

func PublishVoiceEvent(userID, chatID, action string) {
	if kafkaWriter == nil {
		log.Print("дурацкий nil")
		return
	}

	event := VoiceEvent{
		UserId: userID,
		ChatId: chatID,
		Action: action,
	}
	log.Printf("KAFKA: Preparing to publish Action=%s, User=%s\n", action, userID)
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("KAFKA ERROR: Marshal error: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = kafkaWriter.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(chatID),
			Value: bytes,
		},
	)

	if err != nil {
		log.Printf("Failed to write to kafka: %v", err)
	} else {
		log.Printf("Kafka -> Published: User %s %s chat %s", userID, action, chatID)
	}
}
