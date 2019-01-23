package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(pingHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	if rr.Body.String() != "pong" {
		t.Errorf("Body differs. Expected \"%s\", Got: \"%s\"", "pong", rr.Body.String())
	}
}
