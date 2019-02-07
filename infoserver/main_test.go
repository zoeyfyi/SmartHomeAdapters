package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

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

	expectedRequests := []string{"switchserver/123abc"}

	client = testClient(func(req *http.Request) *http.Response {
		if req.URL.String() != expectedRequests[0] {
			t.Errorf("Expected request \"%s\", actual request \"%s\"", expectedRequests[0], req.URL.String())
		}

		// pop first request of slice
		expectedRequests = expectedRequests[1:]

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("{\"IsOn\":false}")),
			Header:     make(http.Header),
		}
	})
	defer func() {
		client = http.DefaultClient
	}()

	responseSubset := "[{\"id\":\"123abc\",\"nickname\":\"testLightbot\",\"robotType\":\"switch\",\"interfaceType\":" +
		"\"toggle\"},{\"id\":\"T2D2\",\"nickname\":\"testThermoBot\",\"robotType\":" +
		"\"thermostat\",\"interfaceType\":\"range\"}]"

	s := newServer(t)
	defer s.Server.Close()

	req, err := http.NewRequest("GET", s.URL+"/robots", nil)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
	rr := httptest.NewRecorder()

	s.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), responseSubset) {
		t.Errorf("Body differs. Expected at least \"%s\", Got: \"%s\"", responseSubset, rr.Body.String())
	}
}

func TestValidRobotId(t *testing.T) {

	expectedRequests := []string{"http://switchserver/123abc"}

	client = testClient(func(req *http.Request) *http.Response {
		if req.URL.String() != expectedRequests[0] {
			t.Errorf("Expected request \"%s\", actual request \"%s\"", expectedRequests[0], req.URL.String())
		}

		// pop first request of slice
		expectedRequests = expectedRequests[1:]

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("{\"IsOn\":false}")),
			Header:     make(http.Header),
		}
	})
	defer func() {
		client = http.DefaultClient
	}()

	responseSubset := "{\"id\":\"123abc\",\"nickname\":\"testLightbot\",\"robotType\":\"switch\",\"interfaceType\":" +
		"\"toggle\"}"

	s := newServer(t)
	defer s.Server.Close()

	req, err := http.NewRequest("GET", s.URL+"/robot/123abc", nil)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
	rr := httptest.NewRecorder()

	s.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), responseSubset) {
		t.Errorf("Body differs. Expected at least \"%s\", Got: \"%s\"", responseSubset, rr.Body.String())
	}
}

func TestInValidRobotId(t *testing.T) {

	response := "\"No robot with that ID\""

	s := newServer(t)
	defer s.Server.Close()

	req, err := http.NewRequest("GET", s.URL+"/robot/definitelynotalegitrobot", nil)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
	rr := httptest.NewRecorder()

	s.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), response) {
		t.Errorf("Body differs. Expected at least \"%s\", Got: \"%s\"", response, rr.Body.String())
	}

}
