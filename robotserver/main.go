package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader upgrades http connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ignore origin
	},
}

// socket to send messages to
var socket *websocket.Conn

// connectHandler establishes the WebSocket with the client
func connectHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to establish WebSocket connection")

	// upgrade request
	if socket != nil {
		log.Println("Closing existing socket")
		socket.Close()
	}

	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade: ", err)
		return
	}

	log.Println("Socket opened")
	socket = newSocket
}

func sendMessage(w http.ResponseWriter, msg string) {
	if socket == nil {
		log.Println("Robot not connected")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// send LED on message to robot
	if err := socket.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Println("Failed to send message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ledOnHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Turning the LED on")
	sendMessage(w, "led on")
}

func ledOffHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Turning the LED off")
	sendMessage(w, "led off")
}

func main() {
	// register routes
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/led/on", ledOnHandler)
	http.HandleFunc("/led/off", ledOffHandler)

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
