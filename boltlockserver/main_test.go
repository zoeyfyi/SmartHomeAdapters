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
	"github.com/mrbenshef/SmartHomeAdapters/boltlockserver/boltlockserver"
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
	resource, err := pool.Run("smarthomeadapters/boltlockdb", "latest", []string{"POSTGRES_PASSWORD=password"})
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
	boltlockserver.RegisterBoltlockServerServer(s, &server{
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

func TestSuccessfullyAddingBoltlock(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := boltlockserver.NewBoltlockServerClient(conn)
	boltlockRobot, err := client.AddBoltlock(ctx, &boltlockserver.AddBoltlockRequest{
		Id:   "123",
		IsOn: false,
	})

	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

	expectedBoltlock := &boltlockserver.Boltlock{
		Id:   "123",
		IsOn: false,
	}

	if !reflect.DeepEqual(boltlockRobot, expectedBoltlock) {
		t.Errorf("Robots differ. Expected: %+v, Got: %+v", expectedBoltlock, boltlockRobot)
	}
}

func TestAddBoltlockAlreadyAdded(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := boltlockserver.NewBoltlockServerClient(conn)

	boltlockRobot, err := client.AddBoltlock(ctx, &boltlockserver.AddBoltlockRequest{
		Id:   "321",
		IsOn: false,
	})
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

	boltlockRobot, err = client.AddBoltlock(ctx, &boltlockserver.AddBoltlockRequest{
		Id:   "321",
		IsOn: false,
	})

	if boltlockRobot != nil {
		t.Errorf("Expected nil boltlock to be returned, got: %+v", boltlockRobot)
	}

	expectedError := "rpc error: code = AlreadyExists desc = Robot \"321\" is already a registered boltlock"
	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, got error: %s", expectedError, err.Error())
	}
}

func TestSuccessfullyRemovingBoltlock(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := boltlockserver.NewBoltlockServerClient(conn)
	_, err = client.RemoveBoltlock(ctx, &boltlockserver.RemoveBoltlockRequest{
		Id: "123",
	})

	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
}

func TestRemoveBoltlockDoesntExist(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := boltlockserver.NewBoltlockServerClient(conn)

	_, err = client.RemoveBoltlock(ctx, &boltlockserver.RemoveBoltlockRequest{
		Id: "doesntexist",
	})

	expectedError := "rpc error: code = InvalidArgument desc = Robot \"doesntexist\" is not a boltlock"
	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, got error: %s", expectedError, err.Error())
	}
}

// TODO: re-add on/off tests once we have calibration routes
