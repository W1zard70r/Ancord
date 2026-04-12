package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pion/turn/v2"
)

type TURNServer struct {
	server *turn.Server
}

func NewTURNServer() (*TURNServer, error) {
	// Получаем настройки из переменных окружения
	publicIP := os.Getenv("TURN_PUBLIC_IP")
	if publicIP == "" {
		return nil, fmt.Errorf("TURN_PUBLIC_IP обязателен для работы сервера")
	}

	portStr := os.Getenv("TURN_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("неверный TURN_PORT: %v", err)
	}

	// Получаем учетные данные
	username := os.Getenv("TURN_USERNAME")
	password := os.Getenv("TURN_PASSWORD")
	if username == "" || password == "" {
		return nil, fmt.Errorf("TURN_USERNAME и TURN_PASSWORD обязательны")
	}

	// Создаем UDP слушатель
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(publicIP),
		Port: port,
	}

	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать UDP слушатель: %v", err)
	}

	// Создаем TURN сервер
	server, err := turn.NewServer(turn.ServerConfig{
		Realm: "ancord-turn",
		AuthHandler: func(username, realm string, srcAddr net.Addr) ([]byte, bool) {
			if username == os.Getenv("TURN_USERNAME") {
				// Создаем ключ на основе username:realm:password
				key := turn.GenerateAuthKey(username, realm, os.Getenv("TURN_PASSWORD"))
				return key, true
			}
			return nil, false
		},
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(publicIP),
					Address:      "0.0.0.0",
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать TURN сервер: %v", err)
	}

	log.Printf("TURN сервер запущен на %s:%d", publicIP, port)

	return &TURNServer{server: server}, nil
}

func (t *TURNServer) Close() error {
	return t.server.Close()
}

// Функция для получения TURN credentials в формате для WebRTC
func GetTURNCredentials() (string, string, string) {
	username := os.Getenv("TURN_USERNAME")
	password := os.Getenv("TURN_PASSWORD")
	publicIP := os.Getenv("TURN_PUBLIC_IP")
	port := os.Getenv("TURN_PORT")

	// Создаем временные credentials (действительны 24 часа)
	timestamp := time.Now().Add(24 * time.Hour).Unix()
	tempUsername := fmt.Sprintf("%d:%s", timestamp, username)

	return fmt.Sprintf("turn:%s:%s", publicIP, port), tempUsername, password
}
