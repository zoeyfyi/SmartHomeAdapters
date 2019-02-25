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

	"github.com/mrbenshef/SmartHomeAdapters/thermostatserver/thermostatserver"

	"github.com/ory/dockertest"

	"google.golang.org/grpc/test/bufconn"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	// connect to docker
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// start thermodb
	resource, err := pool.Run("smarthomeadapters/thermodb", "latest", []string{"POSTGRES_PASSWORD=password"})
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
	thermostatserver.RegisterThermostatServerServer(s, &server{
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

func TestGetThermostat(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := thermostatserver.NewThermostatServerClient(conn)

	thermostat, err := client.GetThermostat(context.Background(), &thermostatserver.ThermostatQuery{
		Id: "qwerty",
	})

	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

	expectedThermostat := &thermostatserver.Thermostat{
		Id:            "qwerty",
		Tempreture:    293,
		MinAngle:      30,
		MaxAngle:      170,
		MinTempreture: 283,
		MaxTempreture: 303,
		IsCalibrated:  true,
	}

	if !reflect.DeepEqual(thermostat, expectedThermostat) {
		t.Errorf("Robots differ. Expected: %+v, Got: %+v", expectedThermostat, thermostat)
	}
}
