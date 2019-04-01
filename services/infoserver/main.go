package main

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/infoserver"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/usecaseserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	robotTypeSwitch     = "switch"
	robotTypeThermostat = "thermostat"
)

type server struct {
	DB            *sql.DB
	UsecaseClient usecaseserver.UsecaseServerClient
}

func (s *server) GetRobot(ctx context.Context, query *infoserver.RobotQuery) (*infoserver.Robot, error) {
	var (
		serial    string
		nickname  string
		robotType string
	)

	log.Printf("getting robot with id: %s (user id: %s)", query.Id, query.UserId)

	// query toggleRobots table for matching robots
	row := s.DB.QueryRow(
		"SELECT serial, nickname, robotType FROM robots WHERE serial = $1 AND registeredUserId = $2",
		query.Id,
		query.UserId,
	)

	err := row.Scan(&serial, &nickname, &robotType)
	if err == sql.ErrNoRows {
		return nil, status.Newf(codes.NotFound, "Robot \"%s\" does not exist", query.Id).Err()
	} else if err != nil {
		log.Printf("Failed to retrive robot %s: %v", query.Id, err)
		return nil, err
	}

	robot := &infoserver.Robot{
		Id:        serial,
		Nickname:  nickname,
		RobotType: robotType,
	}

	// get the status of the robot
	status, err := s.UsecaseClient.GetStatus(ctx, &usecaseserver.GetStatusRequest{
		Robot: &usecaseserver.Robot{
			Id: serial,
		},
		Usecase: robotType,
	})
	if err != nil {
		return nil, err
	}

	// set the robot interface type and status
	switch status := status.Status.(type) {
	case *usecaseserver.Status_ToggleStatus:
		robot.InterfaceType = "toggle"
		robot.RobotStatus = &infoserver.Robot_ToggleStatus{
			ToggleStatus: &infoserver.ToggleStatus{
				Value: status.ToggleStatus.Value,
			},
		}
	case *usecaseserver.Status_RangeStatus:
		robot.InterfaceType = "range"
		robot.RobotStatus = &infoserver.Robot_RangeStatus{
			RangeStatus: &infoserver.RangeStatus{
				Min:     status.RangeStatus.Min,
				Max:     status.RangeStatus.Max,
				Current: status.RangeStatus.Value,
			},
		}
	}

	return robot, nil
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
		if robotType == robotTypeSwitch {
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

	_, err = s.DB.Exec(
		"INSERT INTO robots (serial, nickname, robotType, registeredUserId) VALUES ($1, $2, $3, $4)",
		query.Id,
		query.Nickname,
		query.RobotType,
		query.UserId,
	)

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) ToggleRobot(ctx context.Context, request *infoserver.ToggleRequest) (*empty.Empty, error) {
	log.Printf("toggling robot %s\n", request.Id)

	robot, err := s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
	if err != nil {
		return nil, err
	}

	// set toggle request
	_, err = s.UsecaseClient.Toggle(ctx, &usecaseserver.ToggleRequest{
		NewValue: request.Value,
		Robot: &usecaseserver.Robot{
			Id: request.Id,
		},
		Usecase: robot.RobotType,
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) RangeRobot(ctx context.Context, request *infoserver.RangeRequest) (*empty.Empty, error) {
	log.Printf("setting range robot %s\n", request.Id)

	robot, err := s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
	if err != nil {
		return nil, err
	}

	// send range request
	_, err = s.UsecaseClient.Range(ctx, &usecaseserver.RangeRequest{
		NewValue: request.Value,
		Robot: &usecaseserver.Robot{
			Id: request.Id,
		},
		Usecase: robot.RobotType,
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) CalibrateRobot(
	ctx context.Context,
	request *infoserver.CalibrationRequest,
) (*infoserver.Robot, error) {
	robot, err := s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
	if err != nil {
		return nil, err
	}

	for _, param := range request.Parameters {
		// basic parameter infomation
		request := &usecaseserver.SetCalibrationParameterRequest{
			Robot: &usecaseserver.Robot{
				Id: robot.Id,
			},
			Id:      param.Id,
			Usecase: robot.RobotType,
		}

		// set the value of the parameter
		switch value := param.Value.(type) {
		case *infoserver.CalibrationParameter_BoolValue:
			request.Details = &usecaseserver.SetCalibrationParameterRequest_BoolValue{
				BoolValue: value.BoolValue,
			}
		case *infoserver.CalibrationParameter_IntValue:
			request.Details = &usecaseserver.SetCalibrationParameterRequest_IntValue{
				IntValue: value.IntValue,
			}
		}

		// submit request
		s.UsecaseClient.SetCalibrationParameter(ctx, request)
	}

	return s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
}

func (s *server) GetCalibration(
	ctx context.Context,
	request *infoserver.RobotQuery,
) (*infoserver.CalibrationParameters, error) {
	robot, err := s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.Id,
		UserId: request.UserId,
	})
	if err != nil {
		return nil, err
	}

	stream, err := s.UsecaseClient.GetCalibrationParameters(ctx, &usecaseserver.GetCalibrationParametersRequest{
		Robot: &usecaseserver.Robot{
			Id: robot.Id,
		},
		Usecase: robot.RobotType,
	})
	if err != nil {
		return nil, err
	}

	parameters := &infoserver.CalibrationParameters{}

	for {
		param, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		infoParam := &infoserver.CalibrationParameter{
			Id:   param.Id,
			Name: param.Name,
		}

		switch details := param.Details.(type) {
		case *usecaseserver.CalibrationParameter_BoolParameter:
			infoParam.Value = &infoserver.CalibrationParameter_BoolValue{
				BoolValue: details.BoolParameter.Current,
			}
		case *usecaseserver.CalibrationParameter_IntParameter:
			infoParam.Value = &infoserver.CalibrationParameter_IntValue{
				IntValue: details.IntParameter.Current,
			}
		}

		parameters.Parameters = append(parameters.Parameters, infoParam)
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

func main() {
	db, err := microservice.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer db.Close()

	// test database
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %+v\n", db.Stats())

	// connect to usecaseserver
	usecaseserverConn, err := grpc.Dial("usecaseserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer usecaseserverConn.Close()
	usecaseClient := usecaseserver.NewUsecaseServerClient(usecaseserverConn)

	// start grpc server
	grpcServer := grpc.NewServer()
	infoServer := &server{
		DB:            db,
		UsecaseClient: usecaseClient,
	}
	infoserver.RegisterInfoServerServer(grpcServer, infoServer)
	lis, err := net.Listen("tcp", "infoserver:80")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve GRPC server: %v", err)
	}
}
