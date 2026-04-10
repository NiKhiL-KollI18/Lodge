package lib

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	rooms    = make(map[string]*Room)
	globalMu sync.Mutex
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRoomID(length int) string {
	seedRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seedRand.Intn(len(charset))]
	}
	return string(b)
}

func CreateRoom(capacity int) (*Room, string) {
	globalMu.Lock()
	defer globalMu.Unlock()

	var roomID string
	//collision protection
	for {
		roomID = GenerateRoomID(6)
		if _, exists := rooms[roomID]; !exists {
			break
		}
	}

	room := &Room{
		ID:        roomID,
		Peers:     make(map[*websocket.Conn]bool),
		Capacity:  capacity,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	rooms[roomID] = room
	return room, roomID
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
