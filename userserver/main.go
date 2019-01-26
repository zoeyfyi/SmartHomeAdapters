package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

// error messages for register route
const (
	ErrorInvalidJSON   = "Invalid JSON"
	ErrorEmailBlank    = "Email is blank"
	ErrorEmailInvalid  = "Email address is invalid"
	ErrorPasswordBlank = "Password is blank"
)

type restError struct {
	Error string `json:"error"`
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

type registerBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/register")

	// decode body
	var registerBody registerBody
	err := json.NewDecoder(r.Body).Decode(&registerBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(restError{ErrorInvalidJSON})
		return
	}

	// check fields are correct
	if registerBody.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(restError{ErrorEmailBlank})
		return
	}
	if !isEmailValid(registerBody.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(restError{ErrorEmailInvalid})
		return
	}
	if registerBody.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(restError{ErrorPasswordBlank})
		return
	}
}

func main() {
	// register routes
	http.HandleFunc("/register", registerHandler)

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
