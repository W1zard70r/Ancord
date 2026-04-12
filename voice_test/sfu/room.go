package main

import (
	"log"
	"sync"

	"github.com/pion/webrtc/v3"
)

type Room struct {
	sync.RWMutex
	Peers  []*webrtc.PeerConnection
	Tracks []*webrtc.TrackLocalStaticRTP
}

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) Join(pc *webrtc.PeerConnection) {
	r.Lock()
	defer r.Unlock()
	r.Peers = append(r.Peers, pc)
}

func (r *Room) Leave(pc *webrtc.PeerConnection) {
	r.Lock()
	defer r.Unlock()

	for i, peer := range r.Peers {
		if peer == pc {
			r.Peers = append(r.Peers[:i], r.Peers[i+1:]...)
			return
		}
	}
}

func (r *Room) AddTrack(newTrack *webrtc.TrackLocalStaticRTP, owner *webrtc.PeerConnection) {
	r.Lock()
	defer r.Unlock()

	r.Tracks = append(r.Tracks, newTrack)

	for _, pc := range r.Peers {
		if pc == owner {
			continue
		}

		if _, err := pc.AddTrack(newTrack); err != nil {
			log.Printf("Ошибка добавления нового трека участнику: %v", err)
		}
	}
}

func (r *Room) SyncTracks(pc *webrtc.PeerConnection) {
	r.RLock()
	defer r.RUnlock()

	for _, existingTrack := range r.Tracks {
		if _, err := pc.AddTrack(existingTrack); err != nil {
			log.Printf("Ошибка добавления существующего трека: %v", err)
		}
	}
}
