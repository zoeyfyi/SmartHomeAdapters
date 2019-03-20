package main

import (
	"context"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc/test/bufconn"

	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/switchserver"
	"google.golang.org/grpc"
)

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	username := os.Getenv("DB_USERNAME")
	if username != "temp" {
		log.Fatalf("Database username must be \"temp\", data will be wiped!")
	}

	db, err := microservice.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// start test gRPC server
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	switchserver.RegisterSwitchServerServer(s, &server{
		DB: db,
	})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	exitCode := m.Run()

	os.Exit(exitCode)
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestSuccessfullyAddingSwitch(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := switchserver.NewSwitchServerClient(conn)
	switchRobot, err := client.AddSwitch(ctx, &switchserver.AddSwitchRequest{
		Id:   "123",
		IsOn: false,
	})

	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

	expectedSwitch := &switchserver.Switch{
		Id:   "123",
		IsOn: false,
	}

	if !reflect.DeepEqual(switchRobot, expectedSwitch) {
		t.Errorf("Robots differ. Expected: %+v, Got: %+v", expectedSwitch, switchRobot)
	}
}

func TestAddSwitchAlreadyAdded(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := switchserver.NewSwitchServerClient(conn)

	switchRobot, err := client.AddSwitch(ctx, &switchserver.AddSwitchRequest{
		Id:   "321",
		IsOn: false,
	})
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

	switchRobot, err = client.AddSwitch(ctx, &switchserver.AddSwitchRequest{
		Id:   "321",
		IsOn: false,
	})

	if switchRobot != nil {
		t.Errorf("Expected nil switch to be returned, got: %+v", switchRobot)
	}

	expectedError := "rpc error: code = AlreadyExists desc = Robot \"321\" is already a registered switch"
	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, got error: %s", expectedError, err.Error())
	}
}

func TestSuccessfullyRemovingSwitch(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := switchserver.NewSwitchServerClient(conn)
	_, err = client.RemoveSwitch(ctx, &switchserver.RemoveSwitchRequest{
		Id: "123",
	})

	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
}

func TestRemoveSwitchDoesntExist(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := switchserver.NewSwitchServerClient(conn)

	_, err = client.RemoveSwitch(ctx, &switchserver.RemoveSwitchRequest{
		Id: "doesntexist",
	})

	expectedError := "rpc error: code = InvalidArgument desc = Robot \"doesntexist\" is not a switch"
	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, got error: %s", expectedError, err.Error())
	}
}

// TODO: re-add on/off tests once we have calibration routes
