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
		url = "localhots:5432"
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

func listRobotHandler(w http.ResponseWriter, r *http.Request){
	// Need to make a database connections and retrieve robots
	// TODO - Connect to database here


	resp := &response{
		Id : "123abc",
		Nickname : "myLightRobot",
		Interface: struct {
			RobotType string `json:"type"`
			Min       int `json:"min"`
			Max       int `json:"max"`
		} {RobotType: "lightSwitch", Min: 0, Max: 100} 	}

	jsonResponse, err := json.Marshal(resp)


	if (err != nil){
		w.Write([]byte("Something went wrong with the database. Please try again."))
		panic(err)
	}

	w.Write(jsonResponse)


}

func main() {
	// register routes
	// mux := http.NewServeMux()

	db := getDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %+v\n", db.Stats())
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/robots", listRobotHandler)

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
