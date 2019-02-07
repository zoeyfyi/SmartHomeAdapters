package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var db = getDb()

func clearDatabase(t *testing.T) {
	_, err := db.Exec("DELETE FROM switches WHERE robotId < 9000")
	if err != nil {
		t.Errorf("Error clearing database: %v", err)
	}
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func TestAddSwitchFieldValidation(t *testing.T) {
	clearDatabase(t)

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
	clearDatabase(t)

	req, err := http.NewRequest("POST", "/123", strings.NewReader("{\"isOn\":false}"))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	rr := httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", http.StatusOK, status)
	}

	var robot switchRobot
	err = json.NewDecoder(rr.Body).Decode(&robot)
	if err != nil {
		t.Errorf("Could not read switch robot json: %v", err)
	}

	expectedRobot := switchRobot{
		RobotID: 123,
		IsOn:    false,
	}

	if !reflect.DeepEqual(robot, expectedRobot) {
		t.Errorf("Robots differ. Expected: %+v, Got: %+v", expectedRobot, robot)
	}
}

func TestAddSwitchAlreadyAdded(t *testing.T) {
	clearDatabase(t)

	req, _ := http.NewRequest("POST", "/123", strings.NewReader("{\"isOn\":false}"))
	rr := httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected OK from POST /123")
	}

	req, err := http.NewRequest("POST", "/123", strings.NewReader("{\"isOn\":false}"))
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
	clearDatabase(t)

	req, _ := http.NewRequest("POST", "/123", strings.NewReader("{\"isOn\":false}"))
	rr := httptest.NewRecorder()
	createRouter(db).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected OK from POST /123")
	}

	req, err := http.NewRequest("DELETE", "/123", nil)
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
	clearDatabase(t)

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

func TestTurnSwitchOnOff(t *testing.T) {
	cases := []struct {
		path             string
		expectedRequests []string
	}{
		{
			"/9999/on",
			[]string{"robotserver/servo/90", "robotserver/servo/45"},
		},
		{
			"/9999/off",
			[]string{"robotserver/servo/0", "robotserver/servo/45"},
		},
	}

	for _, c := range cases {
		clearDatabase(t)

		client = testClient(func(req *http.Request) *http.Response {
			if req.URL.String() != c.expectedRequests[0] {
				t.Errorf("Expected request \"%s\", actual request \"%s\"", c.expectedRequests[0], req.URL.String())
			}

			// pop first request of slice
			c.expectedRequests = c.expectedRequests[1:]

			return &http.Response{
				StatusCode: 200,
				Body:       nil,
				Header:     make(http.Header),
			}
		})
		defer func() {
			client = http.DefaultClient
		}()

		req, err := http.NewRequest("PATCH", c.path, nil)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		rr := httptest.NewRecorder()
		createRouter(db).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Status code differs. Expected \"%d\", Got \"%d\" %+v", http.StatusOK, status, rr)
		}
	}

}

func TestTurnSwitchOnOffFailure(t *testing.T) {
	cases := []struct {
		path          string
		expectedCode  int
		expectedError string
	}{
		{"/9998/on", http.StatusBadRequest, ErrorNotCalibrated},
		{"/9997/on", http.StatusBadRequest, ErrorSwitchOn},
		{"/9999/off", http.StatusBadRequest, ErrorSwitchOff},
	}

	for _, c := range cases {
		clearDatabase(t)

		req, err := http.NewRequest("PATCH", c.path, nil)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		rr := httptest.NewRecorder()
		createRouter(db).ServeHTTP(rr, req)

		if status := rr.Code; status != c.expectedCode {
			t.Errorf("Status code differs. Expected \"%d\", Got \"%d\" %+v", c.expectedCode, status, rr)
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
