package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
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
func connectHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		if _, ok := err.(*websocket.CloseError); ok {
			// socket was closed
			log.Println("Robot not connected (socket is closed)")
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			// unknown error
			log.Printf("Failed to send message: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ledOnHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Turning the LED on")
	sendMessage(w, "led on")
}

func ledOffHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Turning the LED off")
	sendMessage(w, "led off")
}

func servoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("Setting servo to %s\n", ps.ByName("angle"))

	angle, err := strconv.Atoi(ps.ByName("angle"))
	if err != nil {
		log.Printf("Angle was not a number")
		return
	}

	if angle > 180 {
		log.Printf("Angle is to large")
		return
	} else if angle < 0 {
		log.Printf("Angle is to small")
		return
	}

	msg := fmt.Sprintf("servo %d", angle)
	sendMessage(w, msg)
}

func createRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/connect", connectHandler)
	router.PUT("/led/on", ledOnHandler)
	router.PUT("/led/off", ledOffHandler)
	router.PUT("/servo/:angle", servoHandler)

	return router
}

func main() {
	// start server
	if err := http.ListenAndServe(":8080", createRouter()); err != nil {
		panic(err)
	}
}
