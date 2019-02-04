package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
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

func TestConnectToWebSocket(t *testing.T) {
	s := newServer(t)
	defer s.Server.Close()

	_, _, err := websocket.DefaultDialer.Dial(s.URL+"/connect", nil)
	if err != nil {
		t.Fatalf("Error dialing: %v", err)
	}
	defer func() { socket = nil }()
}

func TestSendLEDCommand(t *testing.T) {
	requests := []struct {
		path            string
		expectedMessage string
	}{
		{"/led/on", "led on"},
		{"/led/off", "led off"},
	}

	s := newServer(t)
	defer s.Server.Close()

	ws, _, _ := websocket.DefaultDialer.Dial(s.URL+"/connect", nil)
	// NOTE: we need to clear the socket in main.go otherwise it may not close in time
	// before the next test. Once we have add handling for multiple robots we can replace
	// this with `ws.close()`
	defer func() { socket = nil }()

	for _, r := range requests {
		req, err := http.NewRequest("PUT", s.URL+r.path, nil)
		if err != nil {
			t.Errorf("Error with request: %v", err)
		}

		rr := httptest.NewRecorder()
		s.Handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected OK, got: %v", rr.Code)
		}

		typ, msg, err := ws.ReadMessage()
		if err != nil {
			t.Errorf("Error reading websocket message: %v", err)
		}
		if typ != websocket.TextMessage {
			t.Errorf("Expected websocket.TextMessage type, got type: %v", typ)
		}
		if string(msg) != r.expectedMessage {
			t.Errorf("Expected message: \"%s\", got message: \"%s\"", r.expectedMessage, msg)
		}
	}
}

func TestUnavalibleWhenRobotNotConnected(t *testing.T) {
	requests := []string{
		"/led/on",
		"/led/off",
		"/servo/90",
	}

	s := newServer(t)
	defer s.Server.Close()

	for _, r := range requests {
		req, err := http.NewRequest("PUT", s.URL+r, nil)
		if err != nil {
			t.Errorf("Error with request: %v", err)
		}

		rr := httptest.NewRecorder()
		s.Handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected service unavailble, got: %v", rr.Code)
		}

		var restError restError
		err = json.NewDecoder(rr.Body).Decode(&restError)
		if err != nil {
			t.Errorf("Could not read error json: %v", err)
		}

		if restError.Error != ErrNotConnected {
			t.Errorf("Error differs. Expected \"%s\", Got: \"%s\"", ErrNotConnected, restError.Error)
		}
	}
}

func TestSendServoCommand(t *testing.T) {
	s := newServer(t)
	defer s.Server.Close()

	ws, _, _ := websocket.DefaultDialer.Dial(s.URL+"/connect", nil)
	defer func() { socket = nil }()

	req, err := http.NewRequest("PUT", s.URL+"/servo/90", nil)
	if err != nil {
		t.Errorf("Error with request: %v", err)
	}

	rr := httptest.NewRecorder()
	s.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected OK, got: %v", rr.Code)
	}

	typ, msg, err := ws.ReadMessage()
	if err != nil {
		t.Errorf("Error reading websocket message: %v", err)
	}
	if typ != websocket.TextMessage {
		t.Errorf("Expected websocket.TextMessage type, got type: %v", typ)
	}
	if string(msg) != "servo 90" {
		t.Errorf("Expected message: \"%s\", got message: \"%s\"", "servo 90", msg)
	}
}

func TestSendInvalidServoCommand(t *testing.T) {
	requests := []struct {
		path          string
		expectedError string
	}{
		{"/servo/-100", ErrAngleSmall},
		{"/servo/999", ErrAngleLarge},
		{"/servo/foo", ErrAngleNotNumber},
	}

	s := newServer(t)
	defer s.Server.Close()

	for _, r := range requests {
		req, err := http.NewRequest("PUT", s.URL+r.path, nil)
		if err != nil {
			t.Errorf("Error with request: %v", err)
		}

		rr := httptest.NewRecorder()
		s.Handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected bad request, got: %v", rr.Code)
		}

		var restError restError
		err = json.NewDecoder(rr.Body).Decode(&restError)
		if err != nil {
			t.Errorf("Could not read error json: %v", err)
		}

		if restError.Error != r.expectedError {
			t.Errorf("Error differs. Expected \"%s\", Got: \"%s\"", r.expectedError, restError.Error)
		}
	}
}
