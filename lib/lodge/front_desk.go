package lodge

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	rooms    = make(map[string]*Room)
	globalMu sync.Mutex
)

func CreateRoom(roomID string) (*Room, bool) {
	globalMu.Lock()
	defer globalMu.Unlock()

	if _, exists := rooms[roomID]; exists {
		return nil, false
	}

	room := &Room{
		ID:        roomID,
		Peers:     make(map[*websocket.Conn]bool),
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	rooms[roomID] = room
	return room, true
}

func GetRoom(roomID string) (*Room, bool) {
	globalMu.Lock()
	defer globalMu.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		return nil, false
	}
	return room, true
}

func DeleteRoomIfEmpty(roomID string) bool {
	globalMu.Lock()
	defer globalMu.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		return true
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if len(room.Peers) == 0 {
		delete(rooms, roomID)
		log.Println("Retired Room:", roomID)
		return true
	}
	return false
}
