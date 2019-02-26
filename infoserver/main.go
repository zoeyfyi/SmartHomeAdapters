//go:generate protoc --go_out=plugins=grpc:. ./infoserver/infoserver.proto
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
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
	)

	log.Printf("getting robot with id: %s (user id: %s)", query.Id, query.UserId)

	// query toggleRobots table for matching robots
	row := s.DB.QueryRow("SELECT serial, nickname, robotType FROM robots WHERE serial = $1 AND registeredUserId = $2", query.Id, query.UserId)
	err := row.Scan(&serial, &nickname, &robotType)
	if err == sql.ErrNoRows {
		return nil, status.Newf(codes.NotFound, "Robot \"%s\" does not exist", query.Id).Err()
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
	case "thermostat":
		thermostat, err := s.ThermostatClient.GetThermostat(context.Background(), &thermostatserver.ThermostatQuery{
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
			RobotStatus: &infoserver.Robot_RangeStatus{
				RangeStatus: &infoserver.RangeStatus{
					Min:     thermostat.MinTempreture,
					Max:     thermostat.MaxTempreture,
					Current: thermostat.Tempreture,
				},
			},
		}, nil
	default:
		return nil, status.Newf(codes.InvalidArgument, "Invalid robot type \"%s\"", robotType).Err()
	}
}

func (s *server) GetRobots(query *infoserver.RobotsQuery, stream infoserver.InfoServer_GetRobotsServer) error {
	log.Println("getting robots")

	// Query database for robots
	rows, err := s.DB.Query("SELECT serial, nickname, robotType FROM robots WHERE registeredUserId = $1", query.UserId)
	if err != nil {
		log.Printf("Failed to retrive list of robots: %v", err)
		return err
	}

	var (
		serial        string
		nickname      string
		robotType     string
		interfaceType string
	)

	for rows.Next() {
		err := rows.Scan(&serial, &nickname, &robotType)
		if err != nil {
			log.Printf("Failed to scan row of robots table: %v", err)
			return err
		}
		if robotType == "switch" {
			interfaceType = "toggle"
		} else {
			interfaceType = "range"
		}
		err = stream.Send(&infoserver.Robot{
			Id:            serial,
			Nickname:      nickname,
			RobotType:     robotType,
			InterfaceType: interfaceType, // what should this be? used to be "toggle" or "range"
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) RegisterRobot(ctx context.Context, query *infoserver.RegisterRobotQuery) (*empty.Empty, error) {
	log.Println("registering robot")
	rows, err := s.DB.Query("SELECT * FROM robots WHERE serial = $1", query.Id)
	if err != nil {
		log.Println("Failed to search database for robot.")
		return nil, err
	}
	for rows.Next() {
		return nil, status.Newf(codes.AlreadyExists, "Robot \"%s\" already exists", query.Id).Err()
	}

	_, err = s.DB.Exec("INSERT INTO robots (serial, nickname, robotType, registeredUserId) VALUES ($1, $2, $3, $4)", query.Id, query.Nickname, query.RobotType, query.UserId)

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) ToggleRobot(ctx context.Context, request *infoserver.ToggleRequest) (*empty.Empty, error) {
	log.Printf("toggling robot %s\n", request.Id)

	// get robot type
	var robotType string
	row := s.DB.QueryRow("SELECT robotType FROM toggleRobots WHERE serial = $1 AND registeredUserId = $2", request.Id, request.UserId)
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

func (s *server) CalibrateRobot(ctx context.Context, request *infoserver.CalibrationRequest) (*infoserver.Robot, error) {
	robot, err := s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
	if err != nil {
		return nil, err
	}

	switch robot.RobotType {
	case "switch":
		parameters := &switchserver.SwitchCalibrationParameters{
			Id: request.Id,
		}

		// parse parameters
		for _, param := range request.Parameters {
			switch param.Name {
			case "OnAngle":
				angle, err := strconv.ParseInt(param.Value, 10, 64)
				if err != nil {
					log.Printf("Failed to parse OnAngle parameter value: %v", err)
					return nil, status.Errorf(codes.InvalidArgument, "Expected integer for field OnAngle")
				}
				parameters.OnAngle = &wrappers.Int64Value{Value: angle}
			case "OffAngle":
				angle, err := strconv.ParseInt(param.Value, 10, 64)
				if err != nil {
					log.Printf("Failed to parse OffAngle parameter value: %v", err)
					return nil, status.Errorf(codes.InvalidArgument, "Expected integer for field OffAngle")
				}
				parameters.OffAngle = &wrappers.Int64Value{Value: angle}
			case "RestAngle":
				angle, err := strconv.ParseInt(param.Value, 10, 64)
				if err != nil {
					log.Printf("Failed to parse RestAngle parameter value: %v", err)
					return nil, status.Errorf(codes.InvalidArgument, "Expected integer for field RestAngle")
				}
				parameters.RestAngle = &wrappers.Int64Value{Value: angle}
			case "IsCalibrated":
				isCalibrated, err := strconv.ParseBool(param.Value)
				if err != nil {
					log.Printf("Failed to parse IsCalibrated parameter value: %v", err)
					return nil, status.Errorf(codes.InvalidArgument, "Expected boolean for field IsCalibrated")
				}
				parameters.IsCalibrated = &wrappers.BoolValue{Value: isCalibrated}
			default:
				return nil, status.Errorf(codes.InvalidArgument, "\"%s\" is not a parameter", param.Name)
			}
		}

		// send calibration request
		_, err := s.SwitchClient.CalibrateSwitch(ctx, parameters)
		if err != nil {
			return nil, err
		}
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "Robot type \"%s\" is not recognized", robot.RobotType)
	}

	return s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
}

func (s *server) GetCalibration(ctx context.Context, request *infoserver.RobotQuery) (*infoserver.CalibrationParameters, error) {
	robot, err := s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
	if err != nil {
		return nil, err
	}

	parameters := &infoserver.CalibrationParameters{}

	switch robot.RobotType {
	case "switch":
		switchRobot, err := s.SwitchClient.GetSwitch(context.Background(), &switchserver.SwitchQuery{
			Id: robot.Id,
		})
		if err != nil {
			return nil, err
		}

		parameters.Parameters = append(parameters.Parameters, &infoserver.CalibrationParameter{
			Name:  "OnAngle",
			Value: fmt.Sprintf("%d", switchRobot.OnAngle),
		})
		parameters.Parameters = append(parameters.Parameters, &infoserver.CalibrationParameter{
			Name:  "OffAngle",
			Value: fmt.Sprintf("%d", switchRobot.OffAngle),
		})
		parameters.Parameters = append(parameters.Parameters, &infoserver.CalibrationParameter{
			Name:  "RestAngle",
			Value: fmt.Sprintf("%d", switchRobot.RestAngle),
		})
		parameters.Parameters = append(parameters.Parameters, &infoserver.CalibrationParameter{
			Name:  "IsCalibrated",
			Value: fmt.Sprintf("%t", switchRobot.IsCalibrated),
		})
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "Robot type \"%s\" is not recognized", robot.RobotType)
	}

	return parameters, nil
}

func (s *server) SetUsecase(ctx context.Context, request *infoserver.SetUsecaseRequest) (*infoserver.Robot, error) {
	res, err := s.DB.Exec("UPDATE robots SET robotType = $1 WHERE serial = $2", request.Usecase, request.Id)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		return nil, status.Newf(codes.Internal, "Failed to update usecase of robot \"%s\"", request.Id).Err()
	}

	// check 1 row was updated
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Failed to get the amount of rows affected: %v", err)
		return nil, status.Newf(codes.Internal, "Internal error").Err()
	}
	if rowsAffected != 1 {
		log.Printf("Expected to update exactly 1 row, rows updated: %d\n", rowsAffected)
		return nil, status.Newf(codes.Internal, "Internal error").Err()
	}

	return s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
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
