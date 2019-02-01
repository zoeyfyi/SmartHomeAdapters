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

type RobotInterface interface{

}
type TriggerInterface struct {
	InterfaceType string `json:"type"`
	StateOn string `json:"stateOn"`
}
type RangeInterface struct {
	InterfaceType string `json:"type"`
	Min int `json:"min"`
	Max int `json:"max"`
}

type Response struct {
	Id string `json:"id"`
	Nickname string `json:"nickname"`
	RobotType string `json:"robotType"`
	RobotInterface `json:"interface"`

}

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


func getResults(rows *sql.Rows, w http.ResponseWriter, r *http.Request, interfaceType string) ([]*Response){
	var (
		serial    string
		nickname  string
		robotType string

		minimum   int
		maximum   int
		stateOn string
	)
	// iterate through rows

	robots := []*Response{}

	for rows.Next() {

		if interfaceType == "toggle" {

			err := rows.Scan(&serial, &nickname, &robotType, &stateOn)
			if err != nil {
				log.Printf("Failed to scan row of toggle table: %v", err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			robotInterface := TriggerInterface{interfaceType, stateOn}
			resp := &Response{
				Id:        serial,
				Nickname:  nickname,
				RobotType: robotType,
				RobotInterface : robotInterface,
			}
			robots = append(robots, resp)


		}   else if interfaceType == "range" {

			err := rows.Scan(&serial, &nickname, &robotType, &minimum, &maximum)
			if err != nil {
				log.Printf("Failed to scan row of range table: %v", err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			robotInterface := RangeInterface{interfaceType, minimum, maximum}

			resp := &Response{
				Id:        serial,
				Nickname:  nickname,
				RobotType: robotType,
				RobotInterface : robotInterface,
			}
			robots = append(robots, resp)

		}    else{
			log.Println("Incorrect interface type specified")
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}

	}
	return robots

}

func listRobotHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		log.Println("/robots")

		// Query database for robots
		t_rows, err := db.Query("SELECT * FROM toggleRobots")
		if err != nil {
			log.Printf("Failed to retrive list of robots: %v", err)
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		robots := getResults(t_rows, w, r, "toggle")
		log.Printf("Writing list of robots: ")

		r_rows, err := db.Query("SELECT * FROM rangeRobots")
		if err != nil {
			log.Printf("Failed to retrive list of robots: %v", err)
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		rangeRobots := getResults(r_rows, w, r, "range")
		robots = append(robots, rangeRobots...)

		json.NewEncoder(w).Encode(robots)

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
