package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)

type Robot struct {
	ID            string `json:"id"`
	Nickname      string `json:"nickname"`
	RobotType     string `json:"robotType"`
	InterfaceType string `json:"interfaceType"`
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

		robotID := ps.ByName("robotId")

		var (
			serial    string
			nickname  string
			robotType string
			minimum   int
			maximum   int
		)

		w.Header().Set("Content-Type", "application/json")

		log.Println("/robot/" + robotID)

		// query toggleRobots table for matching robots
		rows, err := db.Query("SELECT * FROM toggleRobots WHERE serial = $1", robotID)

		// initialise counter
		found := false
		if err != nil {
			log.Printf("Failed to retrive robot %s: %v", robotID, err)
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		for rows.Next() {
			// read from table and write response
			found = true
			err := rows.Scan(&serial, &nickname, &robotType)
			resp := &Robot{
				ID:            serial,
				Nickname:      nickname,
				RobotType:     robotType,
				InterfaceType: "toggle",
			}
			json.NewEncoder(w).Encode(resp)
			if err != nil {
				log.Printf("Failed to scan robot %s: %v", robotID, err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		// if the robot was not found in toggleRobots, search rangeRobots
		if found == false {
			rows, err := db.Query("SELECT * FROM rangeRobots WHERE serial = $1", robotID)

			if err != nil {
				log.Printf("Failed to retrive robot %s: %v", robotID, err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			for rows.Next() {
				found = true
				err := rows.Scan(&serial, &nickname, &robotType, &minimum, &maximum)
				resp := &Robot{
					ID:            serial,
					Nickname:      nickname,
					RobotType:     robotType,
					InterfaceType: "range",
				}
				json.NewEncoder(w).Encode(resp)
				if err != nil {
					log.Printf("Failed to scan robot %s: %v", robotID, err)
					httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}

		}

		// if no robots were found, return nil
		if found == false {
			json.NewEncoder(w).Encode("No robot with that ID")
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

		robots := []*Robot{}

		for rows.Next() {
			err := rows.Scan(&serial, &nickname, &robotType)
			if err != nil {
				log.Printf("Failed to scan row of toggle table: %v", err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			resp := &Robot{
				ID:            serial,
				Nickname:      nickname,
				RobotType:     robotType,
				InterfaceType: "toggle",
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

			resp := &Robot{
				ID:            serial,
				Nickname:      nickname,
				RobotType:     robotType,
				InterfaceType: "range",
			}
			robots = append(robots, resp)
		}

		json.NewEncoder(w).Encode(robots)

	})
}

// proxy forwards the request to a different url
func proxy(method string, url string, w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(method, url, r.Body)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	w.Write(buf.Bytes())
}

// errors for toggle robot route
const (
	ErrNotTogglable = "This robot is not togglable"
)

func toggleRobotHandler(db *sql.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Println(r.URL.String())
		w.Header().Set("Content-Type", "application/json")

		robotID := ps.ByName("robotId")
		value := ps.ByName("value")

		// get robot type
		var robotType string
		row := db.QueryRow("SELECT robotType FROM toggleRobots WHERE serial = $1", robotID)
		err := row.Scan(&robotType)
		if err != nil {
			log.Printf("Failed to retrive list of robots: %v", err)
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		// forward request to relevent service
		switch robotType {
		case "switch":
			switch value {
			case "true":
				url := fmt.Sprintf("http://switchserver/%s/on", robotID)
				proxy(http.MethodPatch, url, w, r)
				return
			case "false":
				url := fmt.Sprintf("http://switchserver/%s/off", robotID)
				proxy(http.MethodPatch, url, w, r)
				return
			}
		default:
			log.Printf("robot type \"%s\" is not toggelable", robotType)
			httpWriteError(w, ErrNotTogglable, http.StatusBadRequest)
			return
		}

	})
}

func createRouter(db *sql.DB) *httprouter.Router {
	router := httprouter.New()
	router.GET("/robot/:robotId", queryRobotHandler(db))
	router.GET("/robots", listRobotHandler(db))
	router.PATCH("/robot/:robotId/toggle/:value", toggleRobotHandler(db))

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
	if err := http.ListenAndServe(":80", createRouter(db)); err != nil {
		panic(err)
	}
}
