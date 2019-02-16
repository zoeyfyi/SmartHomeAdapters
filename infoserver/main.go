//go:generate protoc --go_out=plugins=grpc:. ./infoserver/infoserver.proto
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/infoserver/infoserver"
	"github.com/mrbenshef/SmartHomeAdapters/switchserver/switchserver"
	"github.com/mrbenshef/SmartHomeAdapters/thermostatserver/thermostatserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)

type server struct {
	DB               *sql.DB
	SwitchClient     switchserver.SwitchServerClient
	ThermostatClient thermostatserver.ThermostatServerClient
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
			return nil, status.Newf(codes.NotFound, "No robot with ID \"%s\"", query.Id).Err()
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
		switchRobot, err := s.SwitchClient.GetSwitch(context.Background(), &switchserver.SwitchQuery{
			Id: serial,
		})
		if err != nil {
			return nil, err
		}

		// return robot with status infomation
		return &infoserver.Robot{
			Id:            serial,
			Nickname:      nickname,
			RobotType:     robotType,
			InterfaceType: "toggle",
			RobotStatus: &infoserver.Robot_ToggleStatus{
				ToggleStatus: &infoserver.ToggleStatus{
					Value: switchRobot.IsOn,
				},
			},
		}, nil
	default:
		return nil, status.Newf(codes.InvalidArgument, "Invalid robot type \"%s\"", robotType).Err()
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
		stream, err := s.SwitchClient.SetSwitch(context.Background(), &switchserver.SetSwitchRequest{
			Id:    request.Id,
			On:    request.Value,
			Force: request.Force,
		})
		if err != nil {
			return nil, err
		}

		for {
			status, err := stream.Recv()
			if err != nil {
				return nil, err
			}

			if status.Status == switchserver.SetSwitchStatus_DONE {
				break
			}
		}

	default:
		log.Printf("robot type \"%s\" is not toggelable", robotType)
		return nil, status.Newf(codes.InvalidArgument, "Robot \"%s\" of type \"%s\" cannot be toggled", request.Id, robotType).Err()
	}

	return &empty.Empty{}, nil
}

func (s *server) RangeRobot(ctx context.Context, request *infoserver.RangeRequest) (*empty.Empty, error) {
	log.Printf("setting range robot %s\n", request.Id)

	// get robot type
	var robotType string
	row := s.DB.QueryRow("SELECT robotType FROM rangeRobots WHERE serial = $1", request.Id)
	err := row.Scan(&robotType)
	if err != nil {
		log.Printf("Failed to retrive the robot: %v", err)
		return nil, err
	}

	// forward request to relevent service
	switch robotType {
	case "thermostat":
		stream, err := s.ThermostatClient.SetThermostat(context.Background(), &thermostatserver.SetThermostatRequest{
			Id:         request.Id,
			Tempreture: request.Value,
			Unit:       "celsius",
		})
		if err != nil {
			return nil, err
		}

		for {
			status, err := stream.Recv()
			if err != nil {
				return nil, err
			}

			if status.Status == thermostatserver.SetThermostatStatus_DONE {
				break
			}
		}

	default:
		log.Printf("robot type \"%s\" is not a range robot", robotType)
		return nil, status.Newf(codes.InvalidArgument, "Robot \"%s\" of type \"%s\" is not a range robot", request.Id, robotType).Err()
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

	// connect to switchserver
	switchserverConn, err := grpc.Dial("switchserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer switchserverConn.Close()
	switchClient := switchserver.NewSwitchServerClient(switchserverConn)

	// connect to thermostatserver
	thermostatserverConn, err := grpc.Dial("thermostatserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer thermostatserverConn.Close()
	thermostatClient := thermostatserver.NewThermostatServerClient(thermostatserverConn)

	// start grpc server
	grpcServer := grpc.NewServer()
	infoServer := &server{
		DB:               db,
		SwitchClient:     switchClient,
		ThermostatClient: thermostatClient,
	}
	infoserver.RegisterInfoServerServer(grpcServer, infoServer)
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
}
