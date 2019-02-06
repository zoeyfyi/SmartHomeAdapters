//go:generate protoc --go_out=plugins=grpc:. ./infoserver/infoserver.proto
package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/infoserver/infoserver"
	"google.golang.org/grpc"
)

var client = http.DefaultClient

var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)

type Status interface {
	status()
}

type ToggleStatus struct {
	Value bool `json:"value"`
}

func (s ToggleStatus) status() {}

type RangeStatus struct {
	Min     int `json:"min"`
	Max     int `json:"max"`
	Current int `json:"current"`
}

func (s RangeStatus) status() {}

type Robot struct {
	ID            string `json:"id"`
	Nickname      string `json:"nickname"`
	RobotType     string `json:"robotType"`
	InterfaceType string `json:"interfaceType"`
	RobotStatus   Status `json:"status,omitempty"`
}

type server struct {
	DB *sql.DB
}

func (s *server) GetRobot(ctx context.Context, query *infoserver.RobotQuery) (*infoserver.Robot, error) {
	var (
		serial    string
		nickname  string
		robotType string
		minimum   int
		maximum   int
	)

	log.Println("getting robot with id: " + query.Id)

	// query toggleRobots table for matching robots
	rows, err := s.DB.Query("SELECT * FROM toggleRobots WHERE serial = $1", query.Id)
	if err != nil {
		log.Printf("Failed to retrive robot %s: %v", query.Id, err)
		return nil, err
	}

	found := false

	for rows.Next() {
		// read from table and write response
		found = true
		err := rows.Scan(&serial, &nickname, &robotType)
		if err != nil {
			log.Printf("Failed to scan robot %s: %v", query.Id, err)
			return nil, err
		}
	}

	// if the robot was not found in toggleRobots, search rangeRobots
	if found == false {
		rows, err := s.DB.Query("SELECT * FROM rangeRobots WHERE serial = $1", query.Id)

		if err != nil {
			log.Printf("Failed to retrive robot %s: %v", query.Id, err)
			return nil, err
		}

		for rows.Next() {
			found = true
			err := rows.Scan(&serial, &nickname, &robotType, &minimum, &maximum)
			if err != nil {
				log.Printf("Failed to scan robot %s: %v", query.Id, err)
				return nil, err
			}
		}

	}

	// if no robots were found, return nil
	if found == false {
		return nil, &infoserver.RobotNotFoundError{ID: query.Id}
	}

	// get the status of the robot
	switch robotType {
	case "switch":
		url := fmt.Sprintf("http://switchserver/%s", query.Id)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error executing request: %v", err)
			return nil, err
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		var info map[string]interface{}
		json.Unmarshal(buf.Bytes(), &info)

		log.Printf("Status: %+v", buf.String())

		isOn, ok := info["isOn"].(bool)
		if !ok {
			return nil, &infoserver.FailedRetreiveStatusError{}
		}

		return &infoserver.Robot{
			Id:            serial,
			Nickname:      nickname,
			RobotType:     robotType,
			InterfaceType: "toggle",
			RobotStatus: &infoserver.Robot_ToggleStatus{
				ToggleStatus: &infoserver.ToggleStatus{
					Value: isOn,
				},
			},
		}, nil
	default:
		return nil, &infoserver.InvalidRobotTypeError{Type: robotType}
	}
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
	router.GET("/robots", listRobotHandler(db))
	router.PATCH("/robot/:robotId/toggle/:value", toggleRobotHandler(db))

	return router
}

func main() {
	db := getDb()
	defer db.Close()

	// test database
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %+v\n", db.Stats())

	grpcServer := grpc.NewServer()

	infoServer := &server{
		DB: db,
	}
	infoserver.RegisterInfoServerServer(grpcServer, infoServer)

	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer.Serve(lis)

	// start server
	if err := http.ListenAndServe(":80", createRouter(db)); err != nil {
		panic(err)
	}
}
