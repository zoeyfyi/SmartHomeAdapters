package main

import (
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

type testHandler struct{ *testing.T }

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connectHandler(w, r)
}

type testServer struct {
	Handler http.Handler
	Server  *httptest.Server
	URL     string
}

func newServer(t *testing.T) *testServer {
	var server testServer
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/led/on", ledOnHandler)
	http.HandleFunc("/led/off", ledOffHandler)
	server.Handler = http.DefaultServeMux
	server.Server = httptest.NewServer(server.Handler)
	server.URL = httpToWsProtocol(server.Server.URL)
	return &server
}

func TestConnectToWebSocket(t *testing.T) {
	s := newServer(t)
	defer s.Server.Close()

	ws, _, err := websocket.DefaultDialer.Dial(s.URL+"/connect", nil)
	if err != nil {
		t.Fatalf("Error dialing: %v", err)
	}
	defer ws.Close()
}
