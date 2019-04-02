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
		serial    string
		nickname  string
		robotType string
	)

	for rows.Next() {
		err := rows.Scan(&serial, &nickname, &robotType)
		if err != nil {
			log.Printf("Failed to scan row of robots table: %v", err)
			return err
		}

		robot := &infoserver.Robot{
			Id:        serial,
			Nickname:  nickname,
			RobotType: robotType,
		}

		// get the status of the robot
		status, err := s.UsecaseClient.GetStatus(context.Background(), &usecaseserver.GetStatusRequest{
			Robot: &usecaseserver.Robot{
				Id: serial,
			},
			Usecase: robotType,
		})
		if err != nil {
			return err
		}
		log.Printf("robot status: %v", status)

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

		// send robot
		err = stream.Send(robot)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) RegisterRobot(ctx context.Context, query *infoserver.RegisterRobotQuery) (*empty.Empty, error) {
	log.Printf("registering robot \"%s\"", query.Id)
	row := s.DB.QueryRow("SELECT serial, nickname, robotType, registeredUserId FROM robots WHERE serial = $1", query.Id)

	// get robot
	var (
		serial           string
		nickname         *string
		robotType        *string
		registeredUserID *string
	)
	err := row.Scan(&serial, &nickname, &robotType, &registeredUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Newf(codes.NotFound, "robot \"%s\" does not exist", query.Id).Err()
		}
		log.Printf("error scanning robots: %v", err)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	// check robot is not registerd
	if registeredUserID != nil {
		log.Printf("robot already registered to: %s", *registeredUserID)
		return nil, status.Newf(codes.FailedPrecondition, "robot \"%s\" has already been registered", query.Id).Err()
	}

	// set the usecase
	_, err = s.UsecaseClient.SetUsecase(ctx, &usecaseserver.SetUsecaseRequest{
		Robot: &usecaseserver.Robot{
			Id: query.Id,
		},
		Usecase: query.RobotType,
	})
	if err != nil {
		return nil, err
	}

	// update robot
	_, err = s.DB.Exec(
		"UPDATE robots SET nickname = $1, robotType = $2, registeredUserId = $3 WHERE serial = $4",
		query.Nickname,
		query.RobotType,
		query.UserId,
		query.Id,
	)
	if err != nil {
		log.Printf("failed to update robots: %v", err)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &empty.Empty{}, nil
}

func (s *server) UnregisterRobot(ctx context.Context, query *infoserver.UnregisterRobotQuery) (*empty.Empty, error) {
	log.Printf("unregistering robot \"%s\"", query.Id)

	// update robot
	_, err := s.DB.Exec(
		"UPDATE robots SET registeredUserId = NULL WHERE registeredUserId = $1 AND serial = $2",
		query.UserId,
		query.Id,
	)
	if err != nil {
		log.Printf("failed to update robots: %v", err)
		if err == sql.ErrNoRows {
			return nil, status.Newf(codes.Internal, "robot \"%s\" not registered to your account", query.Id).Err()
		}
		return nil, status.New(codes.Internal, "internal error").Err()
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
		Id:     request.RobotId,
		UserId: request.UserId,
	})
	if err != nil {
		return nil, err
	}
	log.Printf("found robot: %+v", robot)

	// basic parameter infomation
	usecaseRequest := &usecaseserver.SetCalibrationParameterRequest{
		Robot: &usecaseserver.Robot{
			Id: robot.Id,
		},
		Id:      request.Id,
		Usecase: robot.RobotType,
	}

	// set the value of the parameter
	switch value := request.Value.(type) {
	case *infoserver.CalibrationRequest_BoolValue:
		usecaseRequest.Details = &usecaseserver.SetCalibrationParameterRequest_BoolValue{
			BoolValue: value.BoolValue,
		}
	case *infoserver.CalibrationRequest_IntValue:
		usecaseRequest.Details = &usecaseserver.SetCalibrationParameterRequest_IntValue{
			IntValue: value.IntValue,
		}
	}

	// submit request
	_, err = s.UsecaseClient.SetCalibrationParameter(ctx, usecaseRequest)
	if err != nil {
		return nil, err
	}

	return s.GetRobot(ctx, &infoserver.RobotQuery{
		Id:     request.RobotId,
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

		log.Printf("usecase parameter: %+v", param)

		infoParam := &infoserver.CalibrationParameter{
			Id:          param.Id,
			Name:        param.Name,
			Description: param.Description,
		}

		switch details := param.Details.(type) {
		case *usecaseserver.CalibrationParameter_BoolParameter:
			infoParam.Type = "bool"
			infoParam.Details = &infoserver.CalibrationParameter_BoolDetails{
				BoolDetails: &infoserver.BoolDetails{
					Current: details.BoolParameter.Current,
					Default: details.BoolParameter.Default,
				},
			}
		case *usecaseserver.CalibrationParameter_IntParameter:
			infoParam.Type = "int"
			infoParam.Details = &infoserver.CalibrationParameter_IntDetails{
				IntDetails: &infoserver.IntDetails{
					Current: details.IntParameter.Current,
					Default: details.IntParameter.Default,
					Min:     details.IntParameter.Min,
					Max:     details.IntParameter.Max,
				},
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

func (s *server) RenameRobot(ctx context.Context, request *infoserver.RenameRobotRequest) (*empty.Empty, error) {
	res, err := s.DB.Exec("UPDATE robots SET nickname = $1 WHERE serial = $2", request.NewNickname, request.Id)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		return nil, status.Newf(codes.Internal, "Failed to update nickname of robot \"%s\"", request.Id).Err()
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

	return &empty.Empty{}, nil
}

func (s *server) GetUsecases(_ *empty.Empty, stream infoserver.InfoServer_GetUsecasesServer) error {
	inStream, err := s.UsecaseClient.GetUsecases(context.Background(), &empty.Empty{})
	if err != nil {
		return err
	}

	for {
		usecase, err := inStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = stream.Send(&infoserver.Usecase{
			Id:          usecase.Id,
			Name:        usecase.Name,
			Description: usecase.Description,
		})
		if err != nil {
			return err
		}
	}

	return nil
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
