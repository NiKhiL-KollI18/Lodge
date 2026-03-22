package lodge

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Room structure
type Room struct {
	ID        string
	Peers     map[*websocket.Conn]bool
	mu        sync.Mutex
	CreatedAt time.Time
	LastUsed  time.Time
}

func (r *Room) AddPeer(conn *websocket.Conn) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Peers) >= 2 {
		return false
	}
	r.Peers[conn] = true
	r.LastUsed = time.Now()
	return true
}

func (r *Room) RemovePeer(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Peers, conn)
	r.LastUsed = time.Now()
}

func (r *Room) Broadcast(sender *websocket.Conn, message MakeRequest) {
	r.mu.Lock()

	peers := make([]*websocket.Conn, 0, len(r.Peers))
	for peer := range r.Peers {
		if peer != sender {
			peers = append(peers, peer)
		}
	}
	r.LastUsed = time.Now()
	r.mu.Unlock()

	for _, peer := range peers {
		err := peer.WriteJSON(message)
		if err != nil {
			log.Println("Broadcast Error:", err)
		}
	}
}
