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

	"github.com/golang/protobuf/ptypes/empty"
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
	row := s.DB.QueryRow("SELECT * FROM toggleRobots WHERE serial = $1", query.Id)
	err := row.Scan(&serial, &nickname, &robotType)
	if err == sql.ErrNoRows {
		// not in toggleRobots, try rangeRobots
		row := s.DB.QueryRow("SELECT * FROM rangeRobots WHERE serial = $1", query.Id)
		err := row.Scan(&serial, &nickname, &robotType, &minimum, &maximum)
		if err == sql.ErrNoRows {
			// not there either
			return nil, &infoserver.RobotNotFoundError{ID: query.Id}
		} else if err != nil {
			log.Printf("Failed to retrive robot %s: %v", query.Id, err)
			return nil, err
		}
	} else if err != nil {
		log.Printf("Failed to retrive robot %s: %v", query.Id, err)
		return nil, err
	}

	// get the status of the robot
	switch robotType {
	case "switch":
		// TODO: get url from config (https://github.com/mrbenshef/SmartHomeAdapters/issues/79)
		url := fmt.Sprintf("http://switchserver/%s", query.Id)

		// send request
		req, err := http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error executing request: %v", err)
			return nil, err
		}

		// read responce
		isSuccessfull := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !isSuccessfull {
			// try to parse body
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(resp.Body)
			if err != nil {
				return nil, err
			}

			// return error with message from request
			return nil, &infoserver.StatusRequestFailed{
				Message: buf.String(),
			}
		}

		// read status infomation
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		var info map[string]interface{}
		json.Unmarshal(buf.Bytes(), &info)
		isOn, ok := info["isOn"].(bool)
		if !ok {
			return nil, &infoserver.StatusRequestFailed{}
		}

		// return robot with status infomation
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

func (s *server) GetRobots(_ *empty.Empty, stream infoserver.InfoServer_GetRobotsServer) error {
	log.Println("getting robots")

	// Query database for robots
	rows, err := s.DB.Query("SELECT * FROM toggleRobots")
	if err != nil {
		log.Printf("Failed to retrive list of robots: %v", err)
		return err
	}

	var (
		serial    string
		nickname  string
		robotType string
		minimum   int
		maximum   int
	)

	for rows.Next() {
		err := rows.Scan(&serial, &nickname, &robotType)
		if err != nil {
			log.Printf("Failed to scan row of toggle table: %v", err)
			return err
		}

		err = stream.Send(&infoserver.Robot{
			Id:            serial,
			Nickname:      nickname,
			RobotType:     robotType,
			InterfaceType: "toggle",
		})
		if err != nil {
			return err
		}
	}

	rows, err = s.DB.Query("SELECT * FROM rangeRobots")
	if err != nil {
		log.Printf("Failed to retrive list of robots: %v", err)
		return err
	}

	for rows.Next() {
		err := rows.Scan(&serial, &nickname, &robotType, &minimum, &maximum)
		if err != nil {
			log.Printf("Failed to scan row of range table: %v", err)
			return err
		}

		err = stream.Send(&infoserver.Robot{
			Id:            serial,
			Nickname:      nickname,
			RobotType:     robotType,
			InterfaceType: "range",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ToggleRobot(ctx context.Context, request *infoserver.ToggleRequest) (*empty.Empty, error) {
	log.Printf("toggling robot %s\n", request.Id)

	// get robot type
	var robotType string
	row := s.DB.QueryRow("SELECT robotType FROM toggleRobots WHERE serial = $1", request.Id)
	err := row.Scan(&robotType)
	if err != nil {
		log.Printf("Failed to retrive list of robots: %v", err)
		return nil, err
	}

	// forward request to relevent service
	switch robotType {
	case "switch":
		// TODO: get url from config (https://github.com/mrbenshef/SmartHomeAdapters/issues/79)
		var url string
		if request.Value {
			url = fmt.Sprintf("http://switchserver/%s/on", request.Id)
		} else {
			url = fmt.Sprintf("http://switchserver/%s/off", request.Id)
		}

		// encode force option
		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(struct {
			Force bool `json:"force"`
		}{
			Force: request.Force,
		})
		if err != nil {
			return nil, err
		}

		// send request
		req, err := http.NewRequest(http.MethodPut, url, buffer)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error executing request: %v", err)
			return nil, err
		}

		// read responce
		isSuccessfull := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !isSuccessfull {
			// try to parse body
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(resp.Body)
			if err != nil {
				return nil, err
			}

			// return error with message from request
			return nil, &infoserver.ToggleRequestFailed{
				Message: buf.String(),
			}
		}

	default:
		log.Printf("robot type \"%s\" is not toggelable", robotType)
		return nil, &infoserver.RobotNotTogglableError{
			ID:   request.Id,
			Type: robotType,
		}
	}

	return &empty.Empty{}, nil
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
}
