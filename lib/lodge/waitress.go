package lodge

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MakeRequest struct {
	Type   string `json:"type"`
	RoomID string `json:"room_id"`
	Data   string `json:"data,omitempty"`
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
			room, ok := CreateRoom(msg.RoomID)
			if !ok {
				log.Println("Room already exists", msg.RoomID)
				continue
			}
			room.AddPeer(guest)
			currentRoom = msg.RoomID

			log.Println("Created room :", msg.RoomID)

		case "join":
			room, ok := GetRoom(msg.RoomID)
			if !ok {
				log.Println("Room does not exist : ", msg.RoomID)
				continue
			}
			ok = room.AddPeer(guest)
			if !ok {
				log.Println("Room is full:", msg.RoomID)
				continue
			}

			currentRoom = msg.RoomID
			log.Println("Joined room :", msg.RoomID)

			//ping back Peer A
			room.Broadcast(guest, MakeRequest{Type: "peer_joined", RoomID: msg.RoomID})
		case "offer", "answer", "ice":
			if currentRoom != "" {
				room, ok := GetRoom(currentRoom)
				if !ok {
					log.Println("Room does not exist : ", currentRoom)
					continue
				}
				room.Broadcast(guest, msg)
			}

		case "leave":
			room, ok := GetRoom(currentRoom)
			if ok {
				room.RemovePeer(guest)
				DeleteRoomIfEmpty(currentRoom)
			}
			return
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
