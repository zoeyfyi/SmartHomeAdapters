package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mrbenshef/SmartHomeAdapters/robotserver/robotserver"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// httpToWsProtocal converts a http url to a WebSocket url
func httpToWsProtocol(url string) string {
	return "ws" + strings.TrimPrefix(url, "http")
}

type testServer struct {
	Handler http.Handler
	Server  *httptest.Server
	URL     string
}

func newServer(t *testing.T) *testServer {
	var server testServer
	server.Handler = createRouter()
	server.Server = httptest.NewServer(server.Handler)
	server.URL = httpToWsProtocol(server.Server.URL)
	return &server
}

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	// start test gRPC server
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	robotserver.RegisterRobotServerServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	os.Exit(m.Run())
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestConnectToWebSocket(t *testing.T) {
	s := newServer(t)
	defer s.Server.Close()

	_, _, err := websocket.DefaultDialer.Dial(s.URL+"/connect/testrobot", nil)
	if err != nil {
		t.Fatalf("Error dialing: %v", err)
	}
}

func TestSendLEDCommand(t *testing.T) {
	requests := []struct {
		ledRequest      *robotserver.LEDRequest
		expectedMessage string
	}{
		{&robotserver.LEDRequest{RobotId: "testrobot", On: true}, "led on"},
		{&robotserver.LEDRequest{RobotId: "testrobot", On: false}, "led off"},
	}

	s := newServer(t)
	defer s.Server.Close()

	ws, _, _ := websocket.DefaultDialer.Dial(s.URL+"/connect/testrobot", nil)

	for _, r := range requests {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		client := robotserver.NewRobotServerClient(conn)
		_, err = client.SetLED(context.Background(), r.ledRequest)
		if err != nil {
			t.Fatalf("Error with request: %v", err)
		}

		typ, msg, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("Error reading websocket message: %v", err)
		}
		if typ != websocket.TextMessage {
			t.Fatalf("Expected websocket.TextMessage type, got type: %v", typ)
		}
		if string(msg) != r.expectedMessage {
			t.Fatalf("Expected message: \"%s\", got message: \"%s\"", r.expectedMessage, msg)
		}
	}
}

func TestUnavailableWhenRobotNotConnected(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	expectedError := "rpc error: code = Unavailable desc = Robot \"notconnected\" is not connected"

	client := robotserver.NewRobotServerClient(conn)
	_, err = client.SetLED(context.Background(), &robotserver.LEDRequest{RobotId: "notconnected", On: false})
	if err.Error() != expectedError {
		t.Fatalf("Expected error: \"%s\", got error: %s", expectedError, err.Error())
	}

	_, err = client.SetLED(context.Background(), &robotserver.LEDRequest{RobotId: "notconnected", On: true})
	if err.Error() != expectedError {
		t.Fatalf("Expected error: \"%s\", got error: %s", expectedError, err.Error())
	}

	_, err = client.SetServo(context.Background(), &robotserver.ServoRequest{RobotId: "notconnected", Angle: 0})
	if err.Error() != expectedError {
		t.Fatalf("Expected error: \"%s\", got error: %s", expectedError, err.Error())
	}
}

func TestSendServoCommand(t *testing.T) {
	requests := []struct {
		servoRequest    *robotserver.ServoRequest
		expectedMessage string
	}{
		{&robotserver.ServoRequest{RobotId: "testrobot", Angle: 0}, "servo 0"},
		{&robotserver.ServoRequest{RobotId: "testrobot", Angle: 90}, "servo 90"},
		{&robotserver.ServoRequest{RobotId: "testrobot", Angle: 180}, "servo 180"},
	}

	s := newServer(t)
	defer s.Server.Close()

	ws, _, _ := websocket.DefaultDialer.Dial(s.URL+"/connect/testrobot", nil)

	for _, r := range requests {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		client := robotserver.NewRobotServerClient(conn)
		_, err = client.SetServo(context.Background(), r.servoRequest)
		if err != nil {
			t.Fatalf("Error with request: %v", err)
		}

		typ, msg, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("Error reading websocket message: %v", err)
		}
		if typ != websocket.TextMessage {
			t.Fatalf("Expected websocket.TextMessage type, got type: %v", typ)
		}
		if string(msg) != r.expectedMessage {
			t.Fatalf("Expected message: \"%s\", got message: \"%s\"", r.expectedMessage, msg)
		}
	}
}

func TestSendInvalidServoCommand(t *testing.T) {
	requests := []struct {
		servoRequest  *robotserver.ServoRequest
		expectedError string
	}{
		{
			&robotserver.ServoRequest{Angle: -100},
			"rpc error: code = InvalidArgument desc = -100 is to small, must be >= 0",
		},
		{
			&robotserver.ServoRequest{Angle: 230},
			"rpc error: code = InvalidArgument desc = 230 is to large, must be <= 180",
		},
	}

	s := newServer(t)
	defer s.Server.Close()

	for _, r := range requests {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		client := robotserver.NewRobotServerClient(conn)
		_, err = client.SetServo(context.Background(), r.servoRequest)
		if err.Error() != r.expectedError {
			t.Errorf("Expected error: %s, got error: %v", r.expectedError, err)
		}
	}
}
