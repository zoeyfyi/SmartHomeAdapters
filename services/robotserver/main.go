//go:generate protoc --go_out=plugins=grpc:. ./robotserver/robotserver.proto
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/robotserver"
	"google.golang.org/grpc"
)

// upgrader upgrades http connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ignore origin
	},
}

// sockets to send messages to
var sockets = make(map[string]*websocket.Conn)
var socketMutex = &sync.Mutex{}

type server struct{}

// connectHandler establishes the WebSocket with the client
func connectHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Attempting to establish WebSocket connection")

	id := ps.ByName("id")

	// upgrade request
	socketMutex.Lock()
	if sockets[id] != nil {
		log.Println("Closing existing socket")
		sockets[id].Close()
	}
	socketMutex.Unlock()

	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade: ", err)
		return
	}

	log.Println("Socket opened")
	socketMutex.Lock()
	sockets[id] = newSocket
	socketMutex.Unlock()
}

func (s *server) SetServo(ctx context.Context, request *robotserver.ServoRequest) (*empty.Empty, error) {
	log.Printf("setting servo to %d\n", request.Angle)

	if request.Angle > 180 {
		return nil, status.Newf(codes.InvalidArgument, "%d is to large, must be <= 180", request.Angle).Err()
	} else if request.Angle < 0 {
		return nil, status.Newf(codes.InvalidArgument, "%d is to small, must be >= 0", request.Angle).Err()
	}

	msg := fmt.Sprintf("servo %d", request.Angle)
	if err := sendMessage(request.RobotId, msg); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) SetLED(ctx context.Context, request *robotserver.LEDRequest) (*empty.Empty, error) {
	var message string
	if request.On {
		message = "led on"
	} else {
		message = "led off"
	}

	if err := sendMessage(request.RobotId, message); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func sendMessage(id string, msg string) error {
	socketMutex.Lock()
	socket, ok := sockets[id]
	socketMutex.Unlock()

	if !ok || socket == nil {
		// socket was never opened
		log.Println("Socket never opened")
		return status.Newf(codes.Unavailable, "Robot \"%s\" is not connected", id).Err()
	}

	// send LED on message to robot
	if err := socket.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		if _, ok := err.(*websocket.CloseError); ok {
			// socket was closed
			log.Println("Socket closed")
			return status.New(codes.Unavailable, "Robot not connected").Err()
		} else {
			// unknown error
			log.Printf("Failed to send message: %v", err)
			return status.New(codes.Internal, "Failed to communicate with robot").Err()
		}
	}

	return nil
}

func createRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/connect/:id", connectHandler)
	return router
}

func main() {
	// start grpc server
	grpcServer := grpc.NewServer()
	robotServer := &server{}
	robotserver.RegisterRobotServerServer(grpcServer, robotServer)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting grpc server")

	go func() {
		grpcServer.Serve(lis)
	}()

	log.Println("Started grpc server, starting http server")

	// start REST server
	if err := http.ListenAndServe(":80", createRouter()); err != nil {
		panic(err)
	}
}
