package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var db = getDb()

func TestRegisterFieldValidation(t *testing.T) {
	cases := []struct {
		id             string
		body           string
		expectedStatus int
		expectedError  string
	}{
		{"foo", "{\"isOn\":true}", http.StatusBadRequest, ErrorInvalidRobotID},
		{"123", "", http.StatusBadRequest, ErrorInvalidJSON},
		{"foo", "{}", http.StatusBadRequest, ErrorIsOnMissing},
	}

	for _, c := range cases {
		url := fmt.Sprintf("/%s/register", c.id)
		req, err := http.NewRequest("POST", url, strings.NewReader(c.body))
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		rr := httptest.NewRecorder()
		createRouter(db).ServeHTTP(rr, req)

		if status := rr.Code; status != c.expectedStatus {
			t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", c.expectedStatus, status)
		}

		var restError restError
		err = json.NewDecoder(rr.Body).Decode(&restError)
		if err != nil {
			t.Errorf("Could not read error json: %v", err)
		}

		if restError.Error != c.expectedError {
			t.Errorf("Error differs. Expected \"%s\", Got: \"%s\"", c.expectedError, restError.Error)
		}
	}

}

func TestSuccessfullRegistration(t *testing.T) {
	req, err := http.NewRequest("POST", "/123/register", strings.NewReader("{\"isOn\":false}"))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	rr := httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	var switchRobot switchRobot
	err = json.NewDecoder(rr.Body).Decode(&switchRobot)
	if err != nil {
		t.Errorf("Could not read switch robot json: %v", err)
	}

	if switchRobot.IsOn != false {
		t.Errorf("isOn differs. Expected \"%t\", Got: \"%t\"", false, switchRobot.IsOn)
	}
}
