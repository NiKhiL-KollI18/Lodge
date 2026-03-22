package lodge

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type RoomService struct {
	inactiveTTL time.Duration //inactive time limit
	waitingTTL  time.Duration //waiting time
	interval    time.Duration
}

func DeployRoomService(ttl, wtl, interval time.Duration) *RoomService {
	return &RoomService{
		inactiveTTL: ttl,
		waitingTTL:  wtl,
		interval:    interval,
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
		peerCount := len(room.Peers)

		if peerCount == 0 && idleTime > rs.inactiveTTL {
			log.Println("Cleaning empty room : ", id)
			delete(rooms, id)
			room.mu.Unlock()
			continue
		}

		if peerCount == 1 && idleTime > rs.inactiveTTL {
			log.Println("Cleaning inactive room : ", id)

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

		if peerCount > 1 && idleTime > rs.inactiveTTL {
			log.Println("Cleaning over-time room : ", id)

			for client := range room.Peers {
				client.Close()
			}

			delete(rooms, id)
			room.mu.Unlock()
			continue
		}

		room.mu.Unlock()
	}
}
