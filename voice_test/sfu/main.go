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
	port := os.Getenv("PORT")
	InitKafka()
	turnServer, err := NewTURNServer()
	if err != nil {
		log.Fatalf("Failed to start TURN server: %v", err)
	}
	defer turnServer.Close()

	rm := NewRoomManager()
	api := newWebRTCAPI()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		// TODO: доставать user_id из JWT токена
		userID := r.URL.Query().Get("user_id")
		chatID := r.URL.Query().Get("chat_id")
		if userID == "" {
			userID = "unknown_user"
		}
		if chatID == "" {
			chatID = "global_voice"
		}
		room := rm.GetOrCreateRoom(chatID)

		log.Println("WS ATTEMPT: UserID =", userID, "ChatID =", chatID)
		handleSignaling(conn, api, room, userID, chatID, rm)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SFU is running"))
	})

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
		log.Fatalf("Failed to register default codecs: %v", err)
	}
	return webrtc.NewAPI(webrtc.WithMediaEngine(m))
}
