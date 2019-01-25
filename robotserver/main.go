package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader upgrades http connections to WebSocket
var upgrader = websocket.Upgrader{}

// connectHandler establishes the WebSocket with the client
func connectHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade: ", err)
		return
	}
	defer c.Close()

	for {
		// receive a message
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Failed to read message: ", err)
		}
		log.Println("Received: ", message)

		// echo it back
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("Failed to write message: ", err)
			break
		}
	}
}

func main() {
	// register routes
	http.HandleFunc("/connect", connectHandler)

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
