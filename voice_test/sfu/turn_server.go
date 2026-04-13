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

/* TODO: remove credentials generation and add identify logic for any user */

type TURNServer struct {
	server *turn.Server
}

func NewTURNServer() (*TURNServer, error) {
	publicIP := os.Getenv("TURN_PUBLIC_IP")
	if publicIP == "" {
		return nil, fmt.Errorf("TURN_PUBLIC_IP is required")
	}

	portStr := os.Getenv("TURN_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TURN_PORT: %v", err)
	}

	username := os.Getenv("TURN_USERNAME")
	password := os.Getenv("TURN_PASSWORD")
	if username == "" || password == "" {
		return nil, fmt.Errorf("TURN_USERNAME and TURN_PASSWORD are required")
	}

	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(publicIP),
		Port: port,
	}

	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP listener: %v", err)
	}

	server, err := turn.NewServer(turn.ServerConfig{
		Realm: "ancord-turn",
		AuthHandler: func(username, realm string, srcAddr net.Addr) ([]byte, bool) {
			if username == os.Getenv("TURN_USERNAME") {
				// Generate key based on username:realm:password format
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
		return nil, fmt.Errorf("failed to create TURN server: %v", err)
	}

	log.Printf("TURN server started on %s:%d", publicIP, port)

	return &TURNServer{server: server}, nil
}

func (t *TURNServer) Close() error {
	return t.server.Close()
}

func GetTURNCredentials() (string, string, string) {
	username := os.Getenv("TURN_USERNAME")
	password := os.Getenv("TURN_PASSWORD")
	publicIP := os.Getenv("TURN_PUBLIC_IP")
	port := os.Getenv("TURN_PORT")

	// Generate time-limited credentials (valid for 24 hours)
	timestamp := time.Now().Add(24 * time.Hour).Unix()
	tempUsername := fmt.Sprintf("%d:%s", timestamp, username)

	return fmt.Sprintf("turn:%s:%s", publicIP, port), tempUsername, password
}
