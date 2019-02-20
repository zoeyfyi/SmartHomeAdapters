package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ory/dockertest"

	"google.golang.org/grpc/test/bufconn"

	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/switchserver/switchserver"
	"google.golang.org/grpc"
)

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	// connect to docker
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// start infodb
	resource, err := pool.Run("smarthomeadapters/switchdb", "latest", []string{"POSTGRES_PASSWORD=password"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	url = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))

	// wait till db is up
	if err = pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), database))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// start test gRPC server
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	switchserver.RegisterSwitchServerServer(s, &server{
		DB: getDb(),
	})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	exitCode := m.Run()

	pool.Purge(resource)

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
