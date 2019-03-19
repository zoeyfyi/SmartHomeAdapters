//go:generate protoc --go_out=plugins=grpc:. ./switchserver/switchserver.proto
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/robotserver"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/switchserver"
	"google.golang.org/grpc"
)

var robotserverClient robotserver.RobotServerClient

// database connection infomation
var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)

type server struct {
	DB *sql.DB
}

func (s *server) AddSwitch(ctx context.Context, request *switchserver.AddSwitchRequest) (*switchserver.Switch, error) {
	switchRobot := switchserver.Switch{
		Id:   request.Id,
		IsOn: request.IsOn,
	}

	// insert into database
	_, err := s.DB.Exec(
		"INSERT INTO switches(serial, isOn, onAngle, offAngle, restAngle, isCalibrated) VALUES($1, $2, $3, $4, $5, $6)",
		switchRobot.Id,
		switchRobot.IsOn,
		switchRobot.OnAngle,
		switchRobot.OffAngle,
		switchRobot.RestAngle,
		switchRobot.IsCalibrated,
	)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			// robot is already registered
			return nil, status.Newf(codes.AlreadyExists, "Robot \"%s\" is already a registered switch", switchRobot.Id).Err()
		} else {
			log.Printf("Failed to insert user into database: %v", err)
			return nil, status.Newf(codes.Internal, "Could not register robot \"%s\" as a switch", switchRobot.Id).Err()
		}
	}

	return &switchRobot, nil
}

func (s *server) RemoveSwitch(ctx context.Context, request *switchserver.RemoveSwitchRequest) (*empty.Empty, error) {
	result, err := s.DB.Exec("DELETE FROM switches WHERE serial = $1", request.Id)
	if err != nil {
		return nil, status.Newf(codes.Internal, "Failed to unregister switch \"%s\"", request.Id).Err()
	}

	count, err := result.RowsAffected()
	if err != nil {
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	if count == 0 {
		return nil, status.Newf(codes.InvalidArgument, "Robot \"%s\" is not a switch", request.Id).Err()
	}

	return &empty.Empty{}, nil
}

func (s *server) GetSwitch(ctx context.Context, request *switchserver.SwitchQuery) (*switchserver.Switch, error) {
	var (
		isOn         bool
		isCalibrated bool
		onAngle      int
		offAngle     int
		restAngle    int
	)

	row := s.DB.QueryRow("SELECT isOn, isCalibrated, onAngle, offAngle, restAngle FROM switches WHERE serial = $1", request.Id)
	err := row.Scan(&isOn, &isCalibrated, &onAngle, &offAngle, &restAngle)
	if err != nil {
		log.Printf("Failed to scan database: %v", err)
		return nil, status.Newf(codes.Internal, "Failed to fetch switch \"%s\"", request.Id).Err()
	}

	return &switchserver.Switch{
		Id:           request.Id,
		IsOn:         isOn,
		IsCalibrated: isCalibrated,
		OnAngle:      int64(onAngle),
		OffAngle:     int64(offAngle),
		RestAngle:    int64(restAngle),
	}, nil
}

func (s *server) SetSwitch(request *switchserver.SetSwitchRequest, stream switchserver.SwitchServer_SetSwitchServer) error {
	// get switch
	robotSwitch, err := s.GetSwitch(context.Background(), &switchserver.SwitchQuery{
		Id: request.Id,
	})
	if err != nil {
		return err
	}

	// check calibration
	if !robotSwitch.IsCalibrated {
		return status.New(codes.FailedPrecondition, "Switch is not calibrated").Err()
	}

	// if not force check we are going to a different state
	if robotSwitch.IsOn == request.On && !request.Force {
		if robotSwitch.IsOn {
			return status.New(codes.InvalidArgument, "Switch is already on").Err()
		} else {
			return status.New(codes.InvalidArgument, "Switch is already off").Err()
		}
	}

	stream.Send(&switchserver.SetSwitchStatus{
		Status: switchserver.SetSwitchStatus_SETTING,
	})

	var angle int64
	if request.On {
		angle = robotSwitch.OnAngle
	} else {
		angle = robotSwitch.OffAngle
	}

	_, err = robotserverClient.SetServo(context.Background(), &robotserver.ServoRequest{
		RobotId: request.Id,
		Angle:   angle,
	})
	if err != nil {
		return err
	}

	// TODO: replace waiting with message acknowledment
	stream.Send(&switchserver.SetSwitchStatus{
		Status: switchserver.SetSwitchStatus_WAITING,
	})
	time.Sleep(time.Second * 3)

	// return to rest angle
	stream.Send(&switchserver.SetSwitchStatus{
		Status: switchserver.SetSwitchStatus_RETURNING,
	})
	_, err = robotserverClient.SetServo(context.Background(), &robotserver.ServoRequest{
		RobotId: request.Id,
		Angle:   robotSwitch.RestAngle,
	})
	if err != nil {
		return err
	}

	// done
	stream.Send(&switchserver.SetSwitchStatus{
		Status: switchserver.SetSwitchStatus_DONE,
	})

	// update database
	res, err := s.DB.Exec("UPDATE switches SET isOn = $1 WHERE serial = $2", request.On, request.Id)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		return status.Newf(codes.Internal, "Failed to update state of switch \"%s\"", request.Id).Err()
	}

	// check 1 row was updated
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Failed to get the amount of rows affected: %v", err)
		return status.Newf(codes.Internal, "Internal error").Err()
	}
	if rowsAffected != 1 {
		log.Printf("Expected to update exactly 1 row, rows updated: %d\n", rowsAffected)
		return status.Newf(codes.Internal, "Internal error").Err()
	}

	return nil
}

func (s *server) CalibrateSwitch(ctx context.Context, parameters *switchserver.SwitchCalibrationParameters) (*empty.Empty, error) {
	fields := make([]string, 0, 4)
	values := make([]interface{}, 0, 4)

	if parameters.OnAngle != nil {
		fields = append(fields, "OnAngle")
		values = append(values, parameters.OnAngle.GetValue())
	}
	if parameters.OffAngle != nil {
		fields = append(fields, "OffAngle")
		values = append(values, parameters.OffAngle.GetValue())
	}
	if parameters.RestAngle != nil {
		fields = append(fields, "RestAngle")
		values = append(values, parameters.RestAngle.GetValue())
	}
	if parameters.IsCalibrated != nil {
		fields = append(fields, "IsCalibrated")
		values = append(values, parameters.IsCalibrated.GetValue())
	}

	// build query string
	queryString := "UPDATE switches SET "
	for i, field := range fields {
		queryString += fmt.Sprintf("%s = $%d", field, i+1)
		if i != len(fields)-1 {
			queryString += ", "
		}
	}
	queryString += fmt.Sprintf(" WHERE serial = $%d", len(fields)+1)
	log.Printf("query string: %s, values: %v", queryString, append(values, parameters.Id))

	// update database
	res, err := s.DB.Exec(queryString, append(values, parameters.Id)...)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		return nil, status.Newf(codes.Internal, "Failed to update calibration of switch \"%s\"", parameters.Id).Err()
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
	// connect to database
	db := getDb()
	defer db.Close()

	// test connection
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %+v\n", db.Stats())

	// connect to robotserver
	robotserverConn, err := grpc.Dial("robotserver:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer robotserverConn.Close()
	robotserverClient = robotserver.NewRobotServerClient(robotserverConn)

	// start grpc server
	grpcServer := grpc.NewServer()
	switchServer := &server{DB: db}
	switchserver.RegisterSwitchServerServer(grpcServer, switchServer)
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
}
