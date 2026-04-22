package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

// Peer groups a PeerConnection with its signaling websocket and a mutex
type Peer struct {
	PC   *webrtc.PeerConnection
	Conn *websocket.Conn
	mu   sync.Mutex // protects Conn writes
}

// RoomTrack holds a server-local track, its owner and senders per-peer
type RoomTrack struct {
	Track   *webrtc.TrackLocalStaticRTP
	Owner   *Peer
	Senders map[*Peer]*webrtc.RTPSender
}

type Room struct {
	sync.RWMutex
	Peers  []*Peer
	Tracks []*RoomTrack
}

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) Join(p *Peer) {
	r.Lock()
	defer r.Unlock()
	r.Peers = append(r.Peers, p)
}

func (r *Room) Leave(p *Peer) {
	r.Lock()
	defer r.Unlock()

	// remove peer from list
	for i, peer := range r.Peers {
		if peer == p {
			r.Peers = append(r.Peers[:i], r.Peers[i+1:]...)
			break
		}
	}

	// remove tracks owned by leaving peer
	var remaining []*RoomTrack
	for _, rt := range r.Tracks {
		if rt.Owner == p {
			// remove senders from other peers and request renegotiation
			for otherPeer, sender := range rt.Senders {
				if sender == nil || otherPeer == nil || otherPeer.PC == nil {
					continue
				}
				if err := otherPeer.PC.RemoveTrack(sender); err != nil {
					log.Printf("Failed to remove sender from peer: %v", err)
				}
				if otherPeer.Conn != nil {
					otherPeer.mu.Lock()
					if err := otherPeer.Conn.WriteJSON(SignalMsg{Type: "renegotiate", Payload: ""}); err != nil {
						log.Printf("Failed to send renegotiate to peer: %v", err)
					}
					otherPeer.mu.Unlock()
				}
			}
			continue
		}
		remaining = append(remaining, rt)
	}
	r.Tracks = remaining
}

func (r *Room) AddTrack(newTrack *webrtc.TrackLocalStaticRTP, owner *Peer) {
	r.Lock()
	defer r.Unlock()

	rt := &RoomTrack{Track: newTrack, Owner: owner, Senders: make(map[*Peer]*webrtc.RTPSender)}
	r.Tracks = append(r.Tracks, rt)

	// Add the new track to all other participants and request renegotiation
	for _, peer := range r.Peers {
		if peer == owner {
			continue
		}

		sender, err := peer.PC.AddTrack(newTrack)
		if err != nil {
			log.Printf("Failed to add new track to peer: %v", err)
			rt.Senders[peer] = nil
			continue
		}
		rt.Senders[peer] = sender

		if peer.Conn != nil {
			peer.mu.Lock()
			if err := peer.Conn.WriteJSON(SignalMsg{Type: "renegotiate", Payload: ""}); err != nil {
				log.Printf("Failed to send renegotiate to peer: %v", err)
			}
			peer.mu.Unlock()
		}
	}
}

func (r *Room) SyncTracks(peer *Peer) {
	r.RLock()
	defer r.RUnlock()

	// Push all current room tracks to a newly joined participant
	for _, rt := range r.Tracks {
		if rt.Owner == peer {
			continue
		}
		sender, err := peer.PC.AddTrack(rt.Track)
		if err != nil {
			log.Printf("Failed to sync existing track: %v", err)
			continue
		}
		rt.Senders[peer] = sender
	}
}

func (r *Room) RemoveTrack(track *webrtc.TrackLocalStaticRTP) {
	r.Lock()
	defer r.Unlock()

	var remaining []*RoomTrack
	for _, rt := range r.Tracks {
		if rt.Track == track {
			// remove senders from peers and notify renegotiation
			for otherPeer, sender := range rt.Senders {
				if sender == nil || otherPeer == nil || otherPeer.PC == nil {
					continue
				}
				if err := otherPeer.PC.RemoveTrack(sender); err != nil {
					log.Printf("Failed to remove sender from peer: %v", err)
				}
				if otherPeer.Conn != nil {
					otherPeer.mu.Lock()
					if err := otherPeer.Conn.WriteJSON(SignalMsg{Type: "renegotiate", Payload: ""}); err != nil {
						log.Printf("Failed to send renegotiate to peer: %v", err)
					}
					otherPeer.mu.Unlock()
				}
			}
			continue
		}
		remaining = append(remaining, rt)
	}
	r.Tracks = remaining
}
