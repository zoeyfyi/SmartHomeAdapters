package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)
type testServer struct {
	Handler http.Handler
	Server  *httptest.Server
	URL     string
}
func newServer(t *testing.T) *testServer {
	var server testServer
	db := getDb()
	server.Handler = createRouter(db)
	server.Server = httptest.NewServer(server.Handler)
	server.URL = server.Server.URL
	return &server
}


func TestRobots(t *testing.T) {

	responseSubset := "[{\"id\":\"123abc\",\"nickname\":\"testLightbot\",\"robotType\":\"switch\",\"interface\":" +
		"{\"type\":\"toggle\"}},{\"id\":\"T2D2\",\"nickname\":\"testThermoBot\",\"robotType\":" +
		"\"thermostat\",\"interface\":{\"type\":\"range\",\"min\":0,\"max\":100}}]"

	s := newServer(t)
	defer s.Server.Close()

	req, err := http.NewRequest("GET", s.URL + "/robots", nil)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
	rr := httptest.NewRecorder()

	s.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), responseSubset)	{
		t.Errorf("Body differs. Expected at least \"%s\", Got: \"%s\"", responseSubset, rr.Body.String())
	}
}

func TestValidRobotId(t *testing.T) {

	responseSubset := "{\"id\":\"123abc\",\"nickname\":\"testLightbot\",\"robotType\":\"switch\",\"interface\":" +
		"{\"type\":\"toggle\"}}"

	s := newServer(t)
	defer s.Server.Close()

	req, err := http.NewRequest("GET", s.URL + "/robot/123abc", nil)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
	rr := httptest.NewRecorder()

	s.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), responseSubset)	{
		t.Errorf("Body differs. Expected at least \"%s\", Got: \"%s\"", responseSubset, rr.Body.String())
	}
}

func TestInValidRobotId(t *testing.T) {

	response := "\"No robot with that ID\""

	s := newServer(t)
	defer s.Server.Close()

	req, err := http.NewRequest("GET", s.URL + "/robot/definitelynotalegitrobot", nil)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
	rr := httptest.NewRecorder()

	s.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), response)	{
		t.Errorf("Body differs. Expected at least \"%s\", Got: \"%s\"", response, rr.Body.String())
	}

}