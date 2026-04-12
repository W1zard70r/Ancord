package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	// Получаем порт с значением по умолчанию
	port := os.Getenv("PORT")

	// Запускаем TURN сервер
	turnServer, err := NewTURNServer()
	if err != nil {
		log.Fatalf("Ошибка запуска TURN сервера: %v", err)
	}
	defer turnServer.Close()

	room := NewRoom()
	api := newWebRTCAPI()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Ошибка апгрейда WebSocket: %v", err)
			return
		}

		handleSignaling(conn, api, room)
	})

	// УДАЛИТЬ
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SFU is running"))
	})

	// Добавляем endpoint для получения TURN конфигурации
	http.HandleFunc("/turn-config", func(w http.ResponseWriter, r *http.Request) {
		url, username, password := GetTURNCredentials()
		config := map[string]interface{}{
			"iceServers": []map[string]interface{}{
				{
					"urls":       []string{url},
					"username":   username,
					"credential": password,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		json.NewEncoder(w).Encode(config)
	})

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func newWebRTCAPI() *webrtc.API {
	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		log.Fatalf("Не удалось зарегистрировать кодеки: %v", err)
	}
	return webrtc.NewAPI(webrtc.WithMediaEngine(m))
}
