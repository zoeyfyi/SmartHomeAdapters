package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// error messages for register route
const (
	ErrorInvalidJSON   = "Invalid JSON"
	ErrorEmailBlank    = "Email is blank"
	ErrorEmailInvalid  = "Email address is invalid"
	ErrorPasswordBlank = "Password is blank"
	ErrorEmailExists   = "Email already exists"
)

type restError struct {
	Error string `json:"error"`
}

func httpWriteError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(restError{msg})
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

type registerBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("/register")

		// decode body
		var registerBody registerBody
		err := json.NewDecoder(r.Body).Decode(&registerBody)
		if err != nil {
			httpWriteError(w, ErrorInvalidJSON, http.StatusBadRequest)
			return
		}

		email := registerBody.Email
		password := registerBody.Password

		// check fields are correct
		if email == "" {
			httpWriteError(w, ErrorEmailBlank, http.StatusBadRequest)
			return
		}
		if !isEmailValid(email) {
			httpWriteError(w, ErrorEmailInvalid, http.StatusBadRequest)
			return
		}
		if password == "" {
			httpWriteError(w, ErrorPasswordBlank, http.StatusBadRequest)
			return
		}

		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

		// insert user into database
		_, err = db.Exec("INSERT INTO users(email, password) VALUES($1, $2)", email, hash)
		if err != nil {
			log.Printf("Failed to insert user into database: %v", err)
			if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
				// email already exists
				httpWriteError(w, ErrorEmailExists, http.StatusBadRequest)
			} else {
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
	})
}

func main() {
	log.Println("Server starting")

	// database connection infomation
	var (
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		database = os.Getenv("DB_DATABASE")
		url      = os.Getenv("DB_URL")
	)

	connectionStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, url, database)
	log.Printf("Connecting to database with \"%s\"\n", connectionStr)

	db, err := sql.Open("postgres", connectionStr)
	for i := 1; i <= 5; i++ {
		// try to connect
		if err == nil {
			break
		}
		log.Fatalf("(%d) Failed to connect to postgres: %v", i, err)
		time.Sleep(1000)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Println("Connected to database: %v+", db.Stats())

	// register routes
	http.HandleFunc("/register", registerHandler(db))

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
