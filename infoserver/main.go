package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)


type response struct {
	Id string `json:"id"`
	Nickname string `json:"nickname"`
	Interface struct {
		RobotType string `json:"type"`
		Min int `json:"min"`
		Max int `json:"max"`} `json:"interface"`}

func connectionStr() string {
	if username == "" {
		username = "postgres"
	}
	if password == "" {
		password = "password"
	}
	if url == "" {
		url = "localhost:5432"
	}
	if database == "" {
		database = "postgres"
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, url, database)
}

func getDb() *sql.DB {
	log.Printf("Connecting to database with \"%s\"\n", connectionStr())
	db, err := sql.Open("postgres", connectionStr())
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	return db
}

type restError struct {
	Error string `json:"error"`
}
func httpWriteError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(restError{msg})
}
func listRobotHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


		log.Println("/robots")

		// Query database for robots
		rows, err := db.Query("SELECT * FROM robots")
		if err != nil {
			log.Printf("Failed to retrive list of robots: %v", err)
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		var (
			serial    string
			nickname  string
			robotType string
			minimum   int
			maximum   int
		)
		// iterate through rows
		for rows.Next() {
			err := rows.Scan(&serial, &nickname, &robotType, &minimum, &maximum)
			if err != nil {
				log.Printf("Failed to retrive list of robots: %v", err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			// write JSON response
			resp := &response{
				Id:       serial,
				Nickname: nickname,
				Interface: struct {
					RobotType string `json:"type"`
					Min       int    `json:"min"`
					Max       int    `json:"max"`
				}{RobotType: robotType, Min: minimum, Max: maximum}}

			jsonResponse, err := json.Marshal(resp)

			if err != nil {
				log.Printf("Failed to convert robot list into JSON")
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			w.Write(jsonResponse)
		}

	})
}
func main() {
	// register routes
	// mux := http.NewServeMux()

	db := getDb()
	defer db.Close()

	// test database
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %+v\n", db.Stats())

	// set up handlers
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/robots", listRobotHandler(db))

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
