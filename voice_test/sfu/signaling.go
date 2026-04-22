package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

/* TODO: make normal sync and improve signaling logic */
type SignalMsg struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func createPeerConnection(api *webrtc.API) (*webrtc.PeerConnection, error) {
	// Initialize with TURN server credentials for NAT traversal
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
		log.Printf("Failed to create PeerConnection: %v", err)
		return
	}
	defer pc.Close()

	// wrap PeerConnection with its signaling connection
	peer := &Peer{PC: pc, Conn: conn}
	room.Join(peer)
	defer room.Leave(peer)

	// Provide existing tracks in the room to the new participant
	room.SyncTracks(peer)

	pc.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("New incoming track: %s", remoteTrack.ID())

		// Create a local track to broadcast this stream to other participants
		localTrack, err := webrtc.NewTrackLocalStaticRTP(
			remoteTrack.Codec().RTPCodecCapability,
			"audio",
			remoteTrack.ID(),
		)
		if err != nil {
			log.Printf("Failed to create local track: %v", err)
			return
		}

		room.AddTrack(localTrack, peer)

		// Relay RTP until remote track stops, then remove the local track
		go func() {
			relayRTP(remoteTrack, localTrack)
			room.RemoveTrack(localTrack)
		}()
	})

	var pendingCandidates []webrtc.ICECandidateInit

	for {
		var msg SignalMsg
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket closed: %v", err)
			return
		}

		switch msg.Type {
		case "offer":
			if err := handleOffer(peer, msg.Payload, &pendingCandidates); err != nil {
				log.Printf("Failed to handle offer: %v", err)
			}

		case "candidate":
			if err := handleICECandidate(peer, msg.Payload, &pendingCandidates); err != nil {
				log.Printf("Failed to handle ICE candidate: %v", err)
			}

		default:
			log.Printf("Unknown signal type: %s", msg.Type)
		}
	}
}

func handleOffer(p *Peer, sdp string, pending *[]webrtc.ICECandidateInit) error {
	pc := p.PC
	if err := pc.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sdp}); err != nil {
		return err
	}

	// Process any candidates that arrived before the Offer
	for _, c := range *pending {
		if err := pc.AddICECandidate(c); err != nil {
			log.Printf("Failed to add pending ICE candidate: %v", err)
		}
	}
	*pending = nil

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return err
	}

	// Wait for ICE gathering to complete before sending the Answer
	gatherFinished := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(answer); err != nil {
		return err
	}
	<-gatherFinished

	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Conn.WriteJSON(SignalMsg{Type: "answer", Payload: pc.LocalDescription().SDP})
}

func handleICECandidate(p *Peer, payload string, pending *[]webrtc.ICECandidateInit) error {
	candidate := webrtc.ICECandidateInit{Candidate: payload}

	// If RemoteDescription isn't set, we must buffer the candidate
	if p.PC.RemoteDescription() == nil {
		*pending = append(*pending, candidate)
		return nil
	}

	return p.PC.AddICECandidate(candidate)
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
