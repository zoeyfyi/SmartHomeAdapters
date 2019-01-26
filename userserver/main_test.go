package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterFieldValidation(t *testing.T) {
	cases := []struct {
		body           string
		expectedStatus int
		expectedError  string
	}{
		{"", http.StatusBadRequest, ErrorInvalidJSON},
		{"{\"email\":\"\", \"password\":\"bar\"}", http.StatusBadRequest, ErrorEmailBlank},
		{"{\"email\":\"foo@email.com\", \"password\":\"\"}", http.StatusBadRequest, ErrorPasswordBlank},
		{"{\"email\":\"foo\", \"password\":\"bar\"}", http.StatusBadRequest, ErrorEmailInvalid},
	}

	for _, c := range cases {
		req, err := http.NewRequest("GET", "/register", strings.NewReader(c.body))
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(registerHandler).ServeHTTP(rr, req)

		if status := rr.Code; status != c.expectedStatus {
			t.Errorf("Status code differs. Expected \"%d\", Got \"%d\"", c.expectedStatus, status)
		}

		var restError restError
		err = json.NewDecoder(rr.Body).Decode(&restError)
		if err != nil {
			t.Errorf("Could not read error json: %v", err)
		}

		if restError.Error != c.expectedError {
			t.Errorf("Error differs. Expected \"%s\", Got: \"%s\"", c.expectedError, restError)
		}
	}

}
