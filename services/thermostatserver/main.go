//go:generate protoc --go_out=plugins=grpc:. ./thermostatserver/thermostatserver.proto
package main

import (
	"context"
	"database/sql"
	"log"
	"math"
	"net"
	"time"

	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/robotserver"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/thermostatserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	DB          *sql.DB
	RobotClient robotserver.RobotServerClient
}

func (s *server) GetThermostat(
	ctx context.Context,
	query *thermostatserver.ThermostatQuery,
) (*thermostatserver.Thermostat, error) {
	var (
		tempreture    int64
		minAngle      int64
		maxAngle      int64
		minTempreture int64
		maxTempreture int64
		isCalibrated  bool
	)

	row := s.DB.QueryRow(
		`
		SELECT tempreture, minAngle, maxAngle, minTempreture, maxTempreture, isCalibrated 
		FROM thermostats 
		WHERE serial = $1
		`,
		query.Id,
	)
	err := row.Scan(&tempreture, &minAngle, &maxAngle, &minTempreture, &maxTempreture, &isCalibrated)
	if err != nil {
		log.Printf("Failed to scan database: %v", err)
		return nil, status.Newf(codes.Internal, "Failed to fetch thermostat \"%s\"", query.Id).Err()
	}

	return &thermostatserver.Thermostat{
		Id:            query.Id,
		Tempreture:    tempreture,
		MinAngle:      minAngle,
		MaxAngle:      maxAngle,
		MinTempreture: minTempreture,
		MaxTempreture: maxTempreture,
		IsCalibrated:  isCalibrated,
	}, nil
}

func (s *server) SetThermostat(
	request *thermostatserver.SetThermostatRequest,
	stream thermostatserver.ThermostatServer_SetThermostatServer,
) error {
	// get thermostat
	thermostat, err := s.GetThermostat(context.Background(), &thermostatserver.ThermostatQuery{
		Id: request.Id,
	})
	if err != nil {
		return err
	}

	// check calibration
	if !thermostat.IsCalibrated {
		return status.New(codes.FailedPrecondition, "Thermostat is not calibrated").Err()
	}

	// compute target kelvin
	var targetKelvin int64
	switch request.Unit {
	case "kelvin":
		targetKelvin = request.Tempreture
	case "celsius":
		targetKelvin = int64(float64(request.Tempreture) + 273.15)
	case "fahrenheit":
		targetKelvin = int64((float32(request.Tempreture)-32)*5/9 + 273.15)
	default:
		return status.Newf(
			codes.InvalidArgument,
			"\"%s\" is an unrecognized unit of tempreture",
			request.Unit,
		).Err()
	}

	// check tempreture is within range
	if targetKelvin < thermostat.MinTempreture {
		return status.Newf(
			codes.FailedPrecondition,
			"%d %s is less than the minimum tempreture",
			request.Tempreture,
			request.Unit,
		).Err()
	} else if targetKelvin > thermostat.MaxTempreture {
		return status.Newf(
			codes.FailedPrecondition,
			"%d %s is more than the maximum tempreture",
			request.Tempreture,
			request.Unit,
		).Err()
	}

	// compute angle
	tempretureRatio := float64(targetKelvin-thermostat.MinTempreture) /
		float64(thermostat.MaxTempreture-thermostat.MinTempreture)
	angle := float64(thermostat.MinAngle) + float64(thermostat.MaxAngle-thermostat.MinAngle)*tempretureRatio

	err = stream.Send(&thermostatserver.SetThermostatStatus{
		Status: thermostatserver.SetThermostatStatus_SETTING,
	})
	if err != nil {
		log.Fatalf("Failed to send set thermostat status: %v", err)
		return status.Newf(codes.Internal, "internal error").Err()
	}

	// set servo
	_, err = s.RobotClient.SetServo(context.Background(), &robotserver.ServoRequest{
		RobotId: thermostat.Id,
		Angle:   int64(math.Floor(angle)),
	})
	if err != nil {
		return err
	}

	// TODO: replace waiting with message acknowledment
	err = stream.Send(&thermostatserver.SetThermostatStatus{
		Status: thermostatserver.SetThermostatStatus_WAITING,
	})
	if err != nil {
		log.Fatalf("Failed to send set thermostat status: %v", err)
		return status.Newf(codes.Internal, "internal error").Err()
	}
	time.Sleep(time.Second * 3)

	// done
	err = stream.Send(&thermostatserver.SetThermostatStatus{
		Status: thermostatserver.SetThermostatStatus_DONE,
	})
	if err != nil {
		log.Fatalf("Failed to send set thermostat status: %v", err)
		return status.Newf(codes.Internal, "internal error").Err()
	}

	// update database
	res, err := s.DB.Exec("UPDATE thermostats SET tempreture = $1 WHERE serial = $2", targetKelvin, request.Id)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		return status.Newf(codes.Internal, "Failed to update state of the thermostat \"%s\"", request.Id).Err()
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

func main() {
	// connect to database
	db, err := microservice.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer db.Close()

	// test connection
	err = db.Ping()
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
	robotClient := robotserver.NewRobotServerClient(robotserverConn)

	// start grpc server
	grpcServer := grpc.NewServer()
	thermostatServer := &server{
		DB:          db,
		RobotClient: robotClient,
	}
	thermostatserver.RegisterThermostatServerServer(grpcServer, thermostatServer)
	lis, err := net.Listen("tcp", "thermostatserver:80")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
