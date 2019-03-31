package usecaseservice

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/robotserver"
	usecase_service "github.com/mrbenshef/SmartHomeAdapters/microservice/usecase-service"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var robotserverClient robotserver.RobotServerClient

type parameter interface {
	parameter()
}

type IntParameter struct {
	ID      string
	Name    string
	Type    string
	Min     int64
	Max     int64
	Default int64
	Current int64
}

func (p IntParameter) parameter() {}

type BoolParameter struct {
	ID      string
	Name    string
	Type    string
	Default bool
	Current bool
}

func (p BoolParameter) parameter() {}

type Usecase interface {
	SetUsecase(robotID string, db *sql.DB) error
	GetCalibrationParameters(robotID string, db *sql.DB) ([]parameter, error)
	ResetCalibrationParameter(robotID string, id string, db *sql.DB) error
}

// ToggleUsecase controls how a toggle robot behaves
type ToggleUsecase interface {
	Usecase
	GetStatus(robotID string, db *sql.DB) (bool, error)
	SetRobot(robotID string, value bool, db *sql.DB) error
	SetCalibrationParameter(robotID string, id string, value bool, db *sql.DB) error
}

type RangeUsecase interface {
	Usecase
	GetStatus(robotID string, db *sql.DB) (int64, error)
	SetRobot(robotID string, value int64, db *sql.DB) error
	SetCalibrationParameter(robotID string, id string, value int64, db *sql.DB) error
}

// UsecaseServer serves a usecase
type UsecaseServer struct {
	db          *sql.DB
	usecase     interface{}
	robotClient robotserver.RobotServerClient
}

func ServeRangeUsecase(url string, usecase RangeUsecase) error {
	return serveUsecase(url, usecase)
}

func ServeToggleUsecase(url string, usecase ToggleUsecase) error {
	return serveUsecase(url, usecase)
}

// ServeUsecase serves a usecase
func serveUsecase(url string, usecase interface{}) error {
	db, err := microservice.ConnectToDB()
	if err != nil {
		return err
	}

	// connect to robotserver
	robotserverConn, err := grpc.Dial("robotserver:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer robotserverConn.Close()
	robotClient := robotserver.NewRobotServerClient(robotserverConn)

	// start grpc server
	grpcServer := grpc.NewServer()
	usecaseServer := &UsecaseServer{db, usecase, robotClient}
	usecase_service.RegisterUsecaseServerServer(grpcServer, usecaseServer)

	lis, err := net.Listen("tcp", url)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// SetUsecase is called when a robot has been set to this usecase
func (s *UsecaseServer) SetUsecase(ctx context.Context, robot *usecase_service.Robot) (*empty.Empty, error) {
	err := s.usecase.(Usecase).SetUsecase(robot.Id, s.db)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *UsecaseServer) GetStatus(ctx context.Context, robot *usecase_service.Robot) (*usecase_service.Status, error) {
	switch usecase := s.usecase.(type) {
	case ToggleUsecase:
		isOn, err := usecase.GetStatus(robot.Id, s.db)
		if err != nil {
			return nil, err
		}

		return &usecase_service.Status{
			Value: fmt.Sprintf("%t", isOn),
		}, nil
	case RangeUsecase:
		value, err := usecase.GetStatus(robot.Id, s.db)
		if err != nil {
			return nil, err
		}

		return &usecase_service.Status{
			Value: fmt.Sprintf("%d", value),
		}, nil
	default:
		log.Fatalf("unregonized usecase type: %v", usecase)
		return nil, nil // unreachable
	}
}

func (s *UsecaseServer) SetRobot(ctx context.Context, action *usecase_service.Action) (*empty.Empty, error) {
	switch usecase := s.usecase.(type) {
	case ToggleUsecase:
		isOn, err := strconv.ParseBool(action.NewValue)
		if err != nil {
			log.Printf("error parsing bool: %v", err)
			return nil, status.New(codes.FailedPrecondition, "value should be \"true\" or \"false\"").Err()
		}

		err = usecase.SetRobot(action.Robot.Id, isOn, s.db)
		if err != nil {
			return nil, err
		}

		return &empty.Empty{}, nil

	case RangeUsecase:
		isOn, err := strconv.ParseInt(action.NewValue, 10, 64)
		if err != nil {
			log.Printf("error parsing int: %v", err)
			return nil, status.New(codes.FailedPrecondition, "value should be an integer").Err()
		}

		err = usecase.SetRobot(action.Robot.Id, isOn, s.db)
		if err != nil {
			return nil, err
		}

		return &empty.Empty{}, nil

	default:
		log.Fatalf("unregonized usecase type: %v", usecase)
		return nil, nil // unreachable
	}
}

func (s *UsecaseServer) GetCalibrationParameters(robot *usecase_service.Robot, stream usecase_service.UsecaseServer_GetCalibrationParametersServer) error {
	parameters, err := s.usecase.(Usecase).GetCalibrationParameters(robot.Id, s.db)
	if err != nil {
		return err
	}

	for _, p := range parameters {
		switch p := p.(type) {
		case BoolParameter:
			stream.Send(&usecase_service.CalibrationParameter{
				Id:   p.ID,
				Name: p.Name,
				Type: p.Type,
				Details: &usecase_service.CalibrationParameter_BoolParameter{
					BoolParameter: &usecase_service.BoolParameter{
						Default: p.Default,
						Current: p.Current,
					},
				},
			})
		case IntParameter:
			stream.Send(&usecase_service.CalibrationParameter{
				Id:   p.ID,
				Name: p.Name,
				Type: p.Type,
				Details: &usecase_service.CalibrationParameter_IntParameter{
					IntParameter: &usecase_service.IntParameter{
						Min:     p.Min,
						Max:     p.Max,
						Default: p.Default,
						Current: p.Current,
					},
				},
			})
		}
	}

	return nil
}

func (s *UsecaseServer) SetCalibrationParameter(ctx context.Context, request *usecase_service.SetCalibrationParameterRequest) (*empty.Empty, error) {
	switch usecase := s.usecase.(type) {
	case ToggleUsecase:
		err := usecase.SetCalibrationParameter(
			request.Robot.Id,
			request.Id,
			request.Details.(*usecase_service.SetCalibrationParameterRequest_BoolValue).BoolValue,
			s.db,
		)
		if err != nil {
			return nil, err
		}
		return &empty.Empty{}, nil

	case RangeUsecase:
		err := usecase.SetCalibrationParameter(
			request.Robot.Id,
			request.Id,
			request.Details.(*usecase_service.SetCalibrationParameterRequest_IntValue).IntValue,
			s.db,
		)
		if err != nil {
			return nil, err
		}
		return &empty.Empty{}, nil

	default:
		log.Fatalf("unregonized usecase type: %v", usecase)
		return nil, nil // unreachable
	}
}

func (s *UsecaseServer) ResetCalibrationParameter(ctx context.Context, request *usecase_service.ResetCalibrationParameterRequest) (*empty.Empty, error) {
	err := s.usecase.(Usecase).ResetCalibrationParameter(request.Robot.Id, request.Id, s.db)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
