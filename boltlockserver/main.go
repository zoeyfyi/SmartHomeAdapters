//go:generate protoc --go_out=plugins=grpc:. ./boltlockserver/boltlockserver.proto
package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/boltlockserver/boltlockserver"
	"github.com/mrbenshef/SmartHomeAdapters/robotserver/robotserver"
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

func (s *server) AddBoltlock(ctx context.Context, request *boltlockserver.AddBoltlockRequest) (*boltlockserver.Boltlock, error) {
	boltlockRobot := boltlockserver.Boltlock{
		Id:   request.Id,
		IsOn: request.IsOn,
	}

	// insert into database
	_, err := s.DB.Exec(
		"INSERT INTO boltlocks(serial, isOn, onAngle, offAngle, isCalibrated) VALUES($1, $2, $3, $4, $5)",
		boltlockRobot.Id,
		boltlockRobot.IsOn,
		boltlockRobot.OnAngle,
		boltlockRobot.OffAngle,
		boltlockRobot.IsCalibrated,
	)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			// robot is already registered
			return nil, status.Newf(codes.AlreadyExists, "Robot \"%s\" is already a registered boltlock", boltlockRobot.Id).Err()
		} else {
			log.Printf("Failed to insert user into database: %v", err)
			return nil, status.Newf(codes.Internal, "Could not register robot \"%s\" as a boltlock", boltlockRobot.Id).Err()
		}
	}

	return &boltlockRobot, nil
}

func (s *server) RemoveBoltlock(ctx context.Context, request *boltlockserver.RemoveBoltlockRequest) (*empty.Empty, error) {
	result, err := s.DB.Exec("DELETE FROM boltlocks WHERE serial = $1", request.Id)
	if err != nil {
		return nil, status.Newf(codes.Internal, "Failed to unregister boltlock \"%s\"", request.Id).Err()
	}

	count, err := result.RowsAffected()
	if err != nil {
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	if count == 0 {
		return nil, status.Newf(codes.InvalidArgument, "Robot \"%s\" is not a boltlock", request.Id).Err()
	}

	return &empty.Empty{}, nil
}

func (s *server) GetBoltlock(ctx context.Context, request *boltlockserver.BoltlockQuery) (*boltlockserver.Boltlock, error) {
	var (
		isOn         bool
		isCalibrated bool
		onAngle      int
		offAngle     int
	)

	row := s.DB.QueryRow("SELECT isOn, isCalibrated, onAngle, offAngle, restAngle FROM boltlocks WHERE serial = $1", request.Id)
	err := row.Scan(&isOn, &isCalibrated, &onAngle, &offAngle)
	if err != nil {
		log.Printf("Failed to scan database: %v", err)
		return nil, status.Newf(codes.Internal, "Failed to fetch boltlock \"%s\"", request.Id).Err()
	}

	return &boltlockserver.Boltlock{
		Id:           request.Id,
		IsOn:         isOn,
		IsCalibrated: isCalibrated,
		OnAngle:      int64(onAngle),
		OffAngle:     int64(offAngle),
	}, nil
}

func (s *server) SetBoltlock(request *boltlockserver.SetBoltlockRequest, stream boltlockserver.BoltlockServer_SetBoltlockServer) error {
	// get boltlock
	robotBoltlock, err := s.GetBoltlock(context.Background(), &boltlockserver.BoltlockQuery{
		Id: request.Id,
	})
	if err != nil {
		return err
	}

	// check calibration
	if !robotBoltlock.IsCalibrated {
		return status.New(codes.FailedPrecondition, "Bolt lock is not calibrated").Err()
	}

	// if not force check we are going to a different state
	if robotBoltlock.IsOn == request.On && !request.Force {
		if robotBoltlock.IsOn {
			return status.New(codes.InvalidArgument, "Bolt lock is already on").Err()
		} else {
			return status.New(codes.InvalidArgument, "Bolt lock is already off").Err()
		}
	}

	stream.Send(&boltlockserver.SetBoltlockStatus{
		Status: boltlockserver.SetBoltlockStatus_SETTING,
	})

	var angle int64
	if request.On {
		angle = robotBoltlock.OnAngle
	} else {
		angle = robotBoltlock.OffAngle
	}

	_, err = robotserverClient.SetServo(context.Background(), &robotserver.ServoRequest{
		RobotId: request.Id,
		Angle:   angle,
	})
	if err != nil {
		return err
	}

	// done
	stream.Send(&boltlockserver.SetBoltlockStatus{
		Status: boltlockserver.SetBoltlockStatus_DONE,
	})

	// update database
	res, err := s.DB.Exec("UPDATE boltlocks SET isOn = $1 WHERE serial = $2", request.On, request.Id)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		return status.Newf(codes.Internal, "Failed to update state of boltlock \"%s\"", request.Id).Err()
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

func (s *server) CalibrateBoltlock(ctx context.Context, parameters *boltlockserver.BoltlockCalibrationParameters) (*empty.Empty, error) {
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
	if parameters.IsCalibrated != nil {
		fields = append(fields, "IsCalibrated")
		values = append(values, parameters.IsCalibrated.GetValue())
	}

	// build query string
	queryString := "UPDATE boltlocks SET "
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
		return nil, status.Newf(codes.Internal, "Failed to update calibration of boltlock \"%s\"", parameters.Id).Err()
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
	boltlockServer := &server{DB: db}
	boltlockserver.RegisterBoltlockServerServer(grpcServer, boltlockServer)
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
}
