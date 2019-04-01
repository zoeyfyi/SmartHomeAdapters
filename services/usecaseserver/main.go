package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/robotserver"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/usecaseserver"
	usercaseserver "github.com/mrbenshef/SmartHomeAdapters/microservice/usecaseserver"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errInternal = status.New(codes.Internal, "internal error").Err()

var robotserverClient robotserver.RobotServerClient

type Parameter interface {
	parameter()
}

type IntParameter struct {
	ID      string
	Name    string
	Min     int64
	Max     int64
	Default int64
	Current int64
}

func (p IntParameter) parameter() {}

type BoolParameter struct {
	ID      string
	Name    string
	Default bool
	Current bool
}

func (p BoolParameter) parameter() {}

type UsecaseType int

const (
	ToggleUsecaseType UsecaseType = 1
	RangeUsecaseType  UsecaseType = 2
)

// TODO: seperate into toggle and range
type Usecase interface {
	Name() string
	Description() string
	Type() UsecaseType
	DefaultParameters() []Parameter
	GetParameter(id string) *Parameter
	DefaultToggleStatus() bool
	DefaultRangeStatus() (min int, max int, current int)
	Toggle(value bool, parameters []Parameter, controller RobotController) error
	Range(value int64, parameters []Parameter, controller RobotController) error
}

type RobotController struct {
	robotId     string
	robotClient robotserver.RobotServerClient
}

func (c *RobotController) SetServo(angle int64) error {
	// TODO: context!
	_, err := c.robotClient.SetServo(context.Background(), &robotserver.ServoRequest{
		RobotId: c.robotId,
		Angle:   angle,
	})
	return err
}

func (c *RobotController) SetLED(on bool) error {
	// TODO: context!
	_, err := c.robotClient.SetLED(context.Background(), &robotserver.LEDRequest{
		RobotId: c.robotId,
		On:      on,
	})
	return err
}

// UsecaseServer serves a usecase
type UsecaseServer struct {
	db          *sql.DB
	usecases    map[string]Usecase
	robotClient robotserver.RobotServerClient
}

func (s *UsecaseServer) RegisterUsecase(usecase Usecase) {
	s.usecases[usecase.Name()] = usecase
}

// ServeUsecase serves a usecase
func Serve(url string, usecases map[string]Usecase) error {
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
	usecaseServer := &UsecaseServer{db, usecases, robotClient}
	usercaseserver.RegisterUsecaseServerServer(grpcServer, usecaseServer)

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
func (s *UsecaseServer) SetUsecase(ctx context.Context, request *usercaseserver.SetUsecaseRequest) (*empty.Empty, error) {
	// remove old details
	// not fatal if we can't remove them
	_, err := s.db.Exec("DELETE FROM boolparameter WHERE robotId = $1", request.Robot.Id)
	if err != nil {
		log.Printf("error removing boolparameter's: %v", err)
	}
	_, err = s.db.Exec("DELETE FROM intparameter WHERE robotId = $1", request.Robot.Id)
	if err != nil {
		log.Printf("error removing intparameter's: %v", err)
	}
	_, err = s.db.Exec("DELETE FROM togglestatus WHERE robotId = $1", request.Robot.Id)
	if err != nil {
		log.Printf("error removing togglestatus: %v", err)
	}
	_, err = s.db.Exec("DELETE FROM rangestatus WHERE robotId = $1", request.Robot.Id)
	if err != nil {
		log.Printf("error removing rangestatus: %v", err)
	}

	usecase, ok := s.usecases[request.Usecase]
	if !ok {
		return nil, fmt.Errorf("usecase \"%s\" is unrecognized", request.Usecase)
	}

	// insert defaults for parameters
	for _, p := range usecase.DefaultParameters() {
		switch p := p.(type) {
		case BoolParameter:
			s.db.Exec(
				"INSERT INTO boolparameter (serial, robotId, value) VALUES ($1, $2, $3)",
				p.ID,
				request.Robot.Id,
				p.Default,
			)
		case IntParameter:
			s.db.Exec(
				"INSERT INTO intparameter (serial, robotId, value) VALUES ($1, $2, $3)",
				p.ID,
				request.Robot.Id,
				p.Default,
			)
		default:
			log.Fatalf("unrecognized parameter type (%T): %v", p, p)
		}
	}

	// insert default for status
	switch usecase.Type() {
	case ToggleUsecaseType:
		_, err := s.db.Exec(
			"INSERT INTO togglestatus (robotId, value) VALUES ($1, $2)",
			request.Robot.Id,
			usecase.DefaultToggleStatus(),
		)
		if err != nil {
			log.Printf("failed to insert toggle status: %v", err)
			return nil, errInternal
		}

	case RangeUsecaseType:
		_, _, current := usecase.DefaultRangeStatus()
		_, err := s.db.Exec(
			"INSERT INTO rangestatus (robotId, value) VALUES ($1, $2)",
			request.Robot.Id,
			current,
		)
		if err != nil {
			log.Printf("failed to insert toggle status: %v", err)
			return nil, errInternal
		}
	}

	return &empty.Empty{}, nil
}

func (s *UsecaseServer) GetUsecases(_ *empty.Empty, stream usecaseserver.UsecaseServer_GetUsecasesServer) error {
	for _, usecase := range s.usecases {
		err := stream.Send(&usecaseserver.Usecase{
			Name:        usecase.Name(),
			Description: usecase.Description(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UsecaseServer) GetStatus(ctx context.Context, request *usercaseserver.GetStatusRequest) (*usercaseserver.Status, error) {
	usecase, ok := s.usecases[request.Usecase]
	if !ok {
		return nil, fmt.Errorf("usecase \"%s\" is unrecognized", request.Usecase)
	}

	switch usecase.Type() {
	case ToggleUsecaseType:
		row := s.db.QueryRow("SELECT value FROM togglestatus WHERE robotId = $1", request.Robot.Id)
		var value bool
		err := row.Scan(&value)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, status.Newf(codes.NotFound, "robot \"%s\" does not have a toggle status", request.Robot.Id).Err()
			}
			log.Printf("error scanning togglestatus: %v", err)
			return nil, errInternal
		}

		return &usercaseserver.Status{
			Status: &usercaseserver.Status_ToggleStatus{
				ToggleStatus: &usercaseserver.ToggleStatus{
					Value: value,
				},
			},
		}, nil

	case RangeUsecaseType:
		row := s.db.QueryRow("SELECT value FROM rangestatus WHERE robotId = $2", request.Robot.Id)
		var value int
		err := row.Scan(&value)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, status.Newf(codes.NotFound, "robot \"%s\" does not have a range status", request.Robot.Id).Err()
			}
			log.Printf("error scanning rangestatus: %v", err)
			return nil, errInternal
		}

		return &usercaseserver.Status{
			Status: &usercaseserver.Status_RangeStatus{
				RangeStatus: &usercaseserver.RangeStatus{
					Value: int64(value),
				},
			},
		}, nil
	}

	log.Fatalf("unreachable")
	return nil, nil
}

func (s *UsecaseServer) Toggle(ctx context.Context, action *usercaseserver.ToggleRequest) (*empty.Empty, error) {
	usecase, ok := s.usecases[action.Usecase]
	if !ok {
		return nil, fmt.Errorf("usecase \"%s\" is unrecognized", action.Usecase)
	}

	parameters, err := s.getCalibrationParameters(usecase, action.Robot.Id, s.db)
	if err != nil {
		log.Printf("failed to get parameters: %v", err)
		return nil, errInternal
	}

	err = usecase.Toggle(action.NewValue, parameters, RobotController{
		robotClient: s.robotClient,
		robotId:     action.Robot.Id,
	})
	if err != nil {
		return nil, err
	}

	// update status
	res, err := s.db.Exec(
		"UPDATE togglestatus SET value = $1 WHERE robotId = $2",
		action.NewValue,
		action.Robot.Id,
	)
	if err != nil {
		log.Printf("error updating toggle status: %v", err)
		return nil, errInternal
	}
	rows, err := res.RowsAffected()
	log.Printf("updated toggle status, rows effected: %v, err: %v", rows, err)

	return &empty.Empty{}, nil
}

func (s *UsecaseServer) Range(ctx context.Context, action *usercaseserver.RangeRequest) (*empty.Empty, error) {
	usecase, ok := s.usecases[action.Usecase]
	if !ok {
		return nil, fmt.Errorf("usecase \"%s\" is unrecognized", action.Usecase)
	}

	parameters, err := s.getCalibrationParameters(usecase, action.Robot.Id, s.db)
	if err != nil {
		log.Printf("failed to get parameters: %v", err)
		return nil, errInternal
	}

	err = usecase.Range(action.NewValue, parameters, RobotController{
		robotClient: s.robotClient,
		robotId:     action.Robot.Id,
	})
	if err != nil {
		return nil, err
	}

	// update status
	_, err = s.db.Exec(
		"UPDATE rangestatus SET value = $1 WHERE robotId = $2",
		action.NewValue,
		action.Robot.Id,
	)
	if err != nil {
		log.Printf("error updating range status: %v", err)
		return nil, errInternal
	}

	return &empty.Empty{}, nil
}

func (s *UsecaseServer) getCalibrationParameters(usecase Usecase, robotID string, db *sql.DB) ([]Parameter, error) {
	defaultParams := usecase.DefaultParameters()

	params := make([]Parameter, len(defaultParams))

	for _, p := range defaultParams {
		switch p := p.(type) {
		case BoolParameter:
			row := s.db.QueryRow(
				"SELECT value FROM boolparameter WHERE robotId = $1 AND serial = $2",
				robotID,
				p.ID,
			)
			var value bool
			err := row.Scan(&value)
			if err != nil {
				log.Printf("error scanning bool parameter: %v", err)
				return nil, errInternal
			}

			p.Current = value
			params = append(params, p)
		case IntParameter:
			row := s.db.QueryRow(
				"SELECT value FROM intparameter WHERE robotId = $1 AND serial = $2",
				robotID,
				p.ID,
			)
			var value int64
			err := row.Scan(&value)
			if err != nil {
				log.Printf("error scanning int parameter: %v", err)
				return nil, errInternal
			}

			p.Current = value
			params = append(params, p)

		default:
			log.Printf("Error: parameter is not IntParameter or BoolParameter")
			return nil, nil
		}
	}

	return params, nil
}

func (s *UsecaseServer) GetCalibrationParameters(request *usercaseserver.GetCalibrationParametersRequest, stream usercaseserver.UsecaseServer_GetCalibrationParametersServer) error {
	usecase, ok := s.usecases[request.Usecase]
	if !ok {
		return fmt.Errorf("usecase \"%s\" is unrecognized", request.Usecase)
	}
	log.Printf("getting calibration parameter for usecase: %s and robot with id %s", usecase, request.Robot.Id)
	params, err := s.getCalibrationParameters(usecase, request.Robot.Id, s.db)
	if err != nil {
		log.Printf("Error with getCalibrationParameters %v, %v:", err, errInternal)
		return errInternal
	}

	for _, p := range params {
		switch p := p.(type) {
		case BoolParameter:
			log.Printf("Found boolparameter, sending parameters")
			stream.Send(&usercaseserver.CalibrationParameter{
				Id:   p.ID,
				Name: p.Name,
				Details: &usercaseserver.CalibrationParameter_BoolParameter{
					BoolParameter: &usercaseserver.BoolParameter{
						Default: p.Default,
						Current: p.Current,
					},
				},
			})
		case IntParameter:
			log.Printf("Found intparameter, sending parameters")
			stream.Send(&usercaseserver.CalibrationParameter{
				Id:   p.ID,
				Name: p.Name,
				Details: &usercaseserver.CalibrationParameter_IntParameter{
					IntParameter: &usercaseserver.IntParameter{
						Min:     p.Min,
						Max:     p.Max,
						Default: p.Default,
						Current: p.Current,
					},
				},
			})
		default:
			log.Printf("Error: parameter is not IntParameter or BoolParameter")
			return nil
		}
	}

	return nil
}

func (s *UsecaseServer) SetCalibrationParameter(ctx context.Context, request *usercaseserver.SetCalibrationParameterRequest) (*empty.Empty, error) {
	usecase, ok := s.usecases[request.Usecase]
	if !ok {
		return nil, fmt.Errorf("usecase \"%s\" is unrecognized", request.Usecase)
	}

	parameter := usecase.GetParameter(request.Id)
	if parameter == nil {
		return nil, fmt.Errorf("paramter \"%s\" not found", request.Id)
	}

	switch (*parameter).(type) {
	case BoolParameter:
		_, err := s.db.Exec(
			"UPDATE boolparameter SET value = $1 WHERE robotId = $2",
			request.GetBoolValue(),
			request.Robot.Id,
		)
		if err != nil {
			log.Printf("error updating bool parameter: %v", err)
			return nil, errInternal
		}
	case IntParameter:
		_, err := s.db.Exec(
			"UPDATE intparameter SET value = $1 WHERE robotId = $2",
			request.GetIntValue(),
			request.Robot.Id,
		)
		if err != nil {
			log.Printf("error updating int parameter: %v", err)
			return nil, errInternal
		}
	}

	return &empty.Empty{}, nil
}

func main() {
	Serve("usecaseserver:80", map[string]Usecase{
		"switch":   &Switch{},
		"boltlock": &Boltlock{},
	})
}
