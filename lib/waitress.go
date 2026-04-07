package lib

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var globalMaxCapacity = 100

func UpdateGlobalMaxCapacity(capacity int) bool {
	if capacity < 2 {
		return false
	}
	globalMaxCapacity = capacity

	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MakeRequest struct {
	Type       string `json:"type"`
	RoomID     string `json:"room_id"`
	GuestCount int    `json:"guest_count,omitempty"`
	Data       string `json:"data,omitempty"`
}

func Waitress(w http.ResponseWriter, r *http.Request) {
	guest, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrader Error:", err)
		return
	}
	defer guest.Close()

	var currentRoom string

	for {
		var msg MakeRequest
		err := guest.ReadJSON(&msg)
		if err != nil {
			break
		}

		switch msg.Type {
		case "create":

			capacity := msg.GuestCount
			if capacity <= 0 {
				capacity = globalMaxCapacity
			} else if capacity == 1 {
				return
			} else if capacity > globalMaxCapacity {
				capacity = globalMaxCapacity
			}

			room, generatedID := CreateRoom(capacity)

			room.AddPeer(guest)
			currentRoom = generatedID

			guest.WriteJSON(MakeRequest{
				Type:   "room_created",
				RoomID: generatedID,
			})

			log.Printf("Created room %s with capacity : %d\n", msg.RoomID, capacity)

		case "join":
			room, ok := GetRoom(msg.RoomID)
			if !ok {
				log.Println("Room does not exist : ", msg.RoomID)

				guest.WriteJSON(MakeRequest{
					Type: "error",
					Data: "room not found",
				})

				continue
			}
			ok = room.AddPeer(guest)
			if !ok {
				log.Println("Room is full:", msg.RoomID)
				guest.WriteJSON(MakeRequest{
					Type: "error",
					Data: "room is full",
				})
				continue
			}

			currentRoom = msg.RoomID
			log.Println("Joined room :", msg.RoomID)

			//ping back Peer A
			room.Broadcast(guest, MakeRequest{Type: "peer_joined", RoomID: msg.RoomID})

		case "leave":
			room, ok := GetRoom(currentRoom)
			if ok {
				room.RemovePeer(guest)
				DeleteRoomIfEmpty(currentRoom)
			}
			return

		default:
			if currentRoom != "" {
				room, ok := GetRoom(currentRoom)
				if !ok {
					log.Println("Room does not exist : ", currentRoom)
					continue
				}
				room.Broadcast(guest, msg)
			}
		}

	}
	if currentRoom != "" {
		room, ok := GetRoom(currentRoom)
		if ok {
			room.RemovePeer(guest)
			DeleteRoomIfEmpty(currentRoom)
		}
	}
}
