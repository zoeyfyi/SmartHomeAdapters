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

func TestAddSwitchFieldValidation(t *testing.T) {
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
		url := fmt.Sprintf("/%s", c.id)
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

func TestSuccessfullyAddingSwitch(t *testing.T) {
	req, err := http.NewRequest("POST", "/123", strings.NewReader("{\"isOn\":false}"))
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

func TestAddSwitchAlreadyAdded(t *testing.T) {
	req, _ := http.NewRequest("POST", "/321", strings.NewReader("{\"isOn\":false}"))
	rr := httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)

	req, err := http.NewRequest("POST", "/321", strings.NewReader("{\"isOn\":false}"))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	rr = httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusBadRequest, status)
	}

	var restError restError
	err = json.NewDecoder(rr.Body).Decode(&restError)
	if err != nil {
		t.Errorf("Could not read error json: %v", err)
	}

	if restError.Error != ErrorRobotRegistered {
		t.Errorf("Error differs. Expected \"%s\", Got: \"%s\"", ErrorRobotRegistered, restError)
	}
}

func TestSuccessfullyRemovingSwitch(t *testing.T) {
	req, _ := http.NewRequest("POST", "/321", strings.NewReader("{\"isOn\":false}"))
	rr := httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)

	req, err := http.NewRequest("DELETE", "/321", nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	rr = httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}
}

func TestRemoveSwitchDoesntExist(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/543", nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	rr := httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusBadRequest, status)
	}

	var restError restError
	err = json.NewDecoder(rr.Body).Decode(&restError)
	if err != nil {
		t.Errorf("Could not read error json: %v", err)
	}

	if restError.Error != ErrorRobotNotRegistered {
		t.Errorf("Error differs. Expected \"%s\", Got: \"%s\"", ErrorRobotNotRegistered, restError)
	}
}
