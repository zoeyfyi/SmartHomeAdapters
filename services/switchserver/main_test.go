package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc/test/bufconn"

	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/switchserver"
	"google.golang.org/grpc"
)

var lis *bufconn.Listener
var db *sql.DB

func clearDatabase(t *testing.T) {
	_, err := db.Exec("DELETE FROM switches WHERE serial != '123abc'")
	if err != nil {
		t.Fatalf("Error clearing database: %v", err)
	}
}
func TestMain(m *testing.M) {
	username := os.Getenv("DB_USERNAME")
	if username != "temp" {
		log.Fatalf("Database username must be \"temp\", data will be wiped!")
	}

	var err error
	db, err = microservice.ConnectToDB()
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

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// TODO: re-add on/off tests once we have calibration routes
