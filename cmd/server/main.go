package main

import (
	"LODGE/lib/lodge"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	roomService := lodge.DeployRoomService(60*time.Second, 10*60*time.Second, 30*time.Second)
	roomService.Start()

	http.HandleFunc("/signal", lodge.Waitress)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Lodge opened on port:", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server Error:", err)
	}
}
