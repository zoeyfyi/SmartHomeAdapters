package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)

type RobotInterface interface {
}
type TriggerInterface struct {
	InterfaceType string `json:"type"`
}
type RangeInterface struct {
	InterfaceType string `json:"type"`
	Min           int    `json:"min"`
	Max           int    `json:"max"`
}

type Response struct {
	Id             string `json:"id"`
	Nickname       string `json:"nickname"`
	RobotType      string `json:"robotType"`
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

func queryRobotHandler(db *sql.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// get user-supplied ID parameter

		robotId := ps.ByName("robotId")

		var (
			serial    string
			nickname  string
			robotType string
			minimum   int
			maximum   int
		)

		w.Header().Set("Content-Type", "application/json")

		log.Println("/robot/" + robotId)

		// query toggleRobots table for matching robots
		rows, err := db.Query("SELECT * FROM toggleRobots WHERE serial = $1", robotId)

		// initialise counter
		found := false
		if err != nil {
			log.Printf("Failed to retrive robot %s: %v", robotId, err)
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		for rows.Next() {
			// read from table and write response
			found = true
			err := rows.Scan(&serial, &nickname, &robotType)
			robotInterface := TriggerInterface{"toggle"}
			resp := &Response{
				Id:             serial,
				Nickname:       nickname,
				RobotType:      robotType,
				RobotInterface: robotInterface,
			}
			json.NewEncoder(w).Encode(resp)
			if err != nil {
				log.Printf("Failed to scan robot %s: %v", robotId, err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		// if the robot was not found in toggleRobots, search rangeRobots
		if found == false {
			rows, err := db.Query("SELECT * FROM rangeRobots WHERE serial = $1", robotId)

			if err != nil {
				log.Printf("Failed to retrive robot %s: %v", robotId, err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			for rows.Next() {
				found = true
				err := rows.Scan(&serial, &nickname, &robotType, &minimum, &maximum)
				robotInterface := RangeInterface{"range", minimum, maximum}
				resp := &Response{
					Id:             serial,
					Nickname:       nickname,
					RobotType:      robotType,
					RobotInterface: robotInterface,
				}
				json.NewEncoder(w).Encode(resp)
				if err != nil {
					log.Printf("Failed to scan robot %s: %v", robotId, err)
					httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}

		}

		// if no robots were found, return nil
		if found == false {
			json.NewEncoder(w).Encode(nil)
		}

	})
}

func listRobotHandler(db *sql.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		w.Header().Set("Content-Type", "application/json")

		log.Println("/robots")

		// Query database for robots
		rows, err := db.Query("SELECT * FROM toggleRobots")
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

		robots := []*Response{}

		for rows.Next() {
			err := rows.Scan(&serial, &nickname, &robotType)
			if err != nil {
				log.Printf("Failed to scan row of toggle table: %v", err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			robotInterface := TriggerInterface{"toggle"}
			resp := &Response{
				Id:             serial,
				Nickname:       nickname,
				RobotType:      robotType,
				RobotInterface: robotInterface,
			}
			robots = append(robots, resp)
		}

		rows, err = db.Query("SELECT * FROM rangeRobots")
		if err != nil {
			log.Printf("Failed to retrive list of robots: %v", err)
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		for rows.Next() {
			err := rows.Scan(&serial, &nickname, &robotType, &minimum, &maximum)
			if err != nil {
				log.Printf("Failed to scan row of range table: %v", err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			robotInterface := RangeInterface{"range", minimum, maximum}

			resp := &Response{
				Id:             serial,
				Nickname:       nickname,
				RobotType:      robotType,
				RobotInterface: robotInterface,
			}
			robots = append(robots, resp)
		}

		json.NewEncoder(w).Encode(robots)

	})
}

func createRouter(db *sql.DB) *httprouter.Router {
	router := httprouter.New()
	router.GET("/robot/:robotId", queryRobotHandler(db))
	router.GET("/robots", listRobotHandler(db))

	return router
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

	// start server
	if err := http.ListenAndServe(":8080", createRouter(db)); err != nil {
		panic(err)
	}
}
