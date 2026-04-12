package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type SignalMsg struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func createPeerConnection(api *webrtc.API) (*webrtc.PeerConnection, error) {
	// TURN сервер
	turnURL, username, password := GetTURNCredentials()

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs:       []string{turnURL},
				Username:   username,
				Credential: password,
			},
		},
	}

	return api.NewPeerConnection(config)
}

func handleSignaling(conn *websocket.Conn, api *webrtc.API, room *Room) {
	defer conn.Close()

	pc, err := createPeerConnection(api)
	if err != nil {
		log.Printf("Ошибка создания PeerConnection: %v", err)
		return
	}
	defer pc.Close()

	room.Join(pc)
	defer room.Leave(pc)

	room.SyncTracks(pc)

	pc.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Получен входящий поток: %s", remoteTrack.ID())

		localTrack, err := webrtc.NewTrackLocalStaticRTP(
			remoteTrack.Codec().RTPCodecCapability,
			"audio",
			remoteTrack.ID(),
		)
		if err != nil {
			log.Printf("Ошибка создания локального трека: %v", err)
			return
		}

		room.AddTrack(localTrack, pc)

		go relayRTP(remoteTrack, localTrack)
	})

	var pendingCandidates []webrtc.ICECandidateInit

	for {
		var msg SignalMsg
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Соединение WebSocket закрыто: %v", err)
			return
		}

		switch msg.Type {
		case "offer":
			if err := handleOffer(pc, msg.Payload, conn, &pendingCandidates); err != nil {
				log.Printf("Ошибка обработки offer: %v", err)
			}

		case "candidate":
			if err := handleICECandidate(pc, msg.Payload, &pendingCandidates); err != nil {
				log.Printf("Ошибка обработки candidate: %v", err)
			}

		default:
			log.Printf("Неизвестный тип сигнала: %s", msg.Type)
		}
	}
}

func handleOffer(pc *webrtc.PeerConnection, sdp string, conn *websocket.Conn, pending *[]webrtc.ICECandidateInit) error {
	if err := pc.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sdp}); err != nil {
		return err
	}

	for _, c := range *pending {
		if err := pc.AddICECandidate(c); err != nil {
			log.Printf("Ошибка добавления отложенного candidate: %v", err)
		}
	}
	*pending = nil

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return err
	}

	gatherFinished := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(answer); err != nil {
		return err
	}
	<-gatherFinished

	return conn.WriteJSON(SignalMsg{Type: "answer", Payload: pc.LocalDescription().SDP})
}

func handleICECandidate(pc *webrtc.PeerConnection, payload string, pending *[]webrtc.ICECandidateInit) error {
	candidate := webrtc.ICECandidateInit{Candidate: payload}

	if pc.RemoteDescription() == nil {
		*pending = append(*pending, candidate)
		return nil
	}

	return pc.AddICECandidate(candidate)
}

func relayRTP(remoteTrack *webrtc.TrackRemote, localTrack *webrtc.TrackLocalStaticRTP) {
	rtpBuf := make([]byte, 1500)
	for {
		n, _, err := remoteTrack.Read(rtpBuf)
		if err != nil {
			return
		}

		if _, err := localTrack.Write(rtpBuf[:n]); err != nil {
			return
		}
	}
}
