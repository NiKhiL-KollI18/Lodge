package lib

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type RoomService struct {
	inactiveTTL  time.Duration //inactive time limit
	waitingTTL   time.Duration //waiting time
	maxLifeCycle time.Duration //Hard Limit from "createdAt"
	interval     time.Duration
}

func DeployRoomService(ttl, wtl, maxLifetime, interval time.Duration) *RoomService {
	return &RoomService{
		inactiveTTL:  ttl,
		waitingTTL:   wtl,
		maxLifeCycle: maxLifetime,
		interval:     interval,
	}
}

func (rs *RoomService) Start() {
	ticker := time.NewTicker(rs.interval)

	go func() {
		for range ticker.C {
			rs.CleanUp()
		}
	}()
}

func (rs *RoomService) CleanUp() {
	now := time.Now()

	globalMu.Lock()
	defer globalMu.Unlock()

	for id, room := range rooms {

		room.mu.Lock()

		idleTime := now.Sub(room.LastUsed)
		totalAge := now.Sub(room.CreatedAt)
		peerCount := len(room.Peers)

		shouldDelete := false
		reason := ""

		if peerCount == 0 {
			shouldDelete = true
			reason = "empty room"
		}

		if !shouldDelete && rs.maxLifeCycle > 0 && totalAge > rs.maxLifeCycle {
			shouldDelete = true
			reason = "Max Life Reached"
		}

		if !shouldDelete && peerCount == 1 && rs.waitingTTL > 0 && idleTime > rs.waitingTTL {
			shouldDelete = true
			reason = "Waiting Timeout"
		}

		if !shouldDelete && peerCount > 1 && rs.inactiveTTL > 0 && idleTime > rs.inactiveTTL {
			shouldDelete = true
			reason = "Inactive"
		}

		if shouldDelete {
			log.Printf("RoomService : Cleaning Room [%s] | Reason : %s\n ", id, reason)

			var clients []*websocket.Conn
			for client := range room.Peers {
				clients = append(clients, client)
			}

			delete(rooms, id)
			room.mu.Unlock()

			for _, c := range clients {
				c.Close()
			}
			continue
		}
		room.mu.Unlock()
	}
}
