//go:generate protoc --go_out=plugins=grpc:. ./thermostatserver/thermostatserver.proto
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/robotserver/robotserver"
	"github.com/mrbenshef/SmartHomeAdapters/thermostatserver/thermostatserver"
	"google.golang.org/grpc"
)

var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)

type server struct {
	DB          *sql.DB
	RobotClient robotserver.RobotServerClient
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
	robotClient := robotserver.NewRobotServerClient(robotserverConn)

	// start grpc server
	grpcServer := grpc.NewServer()
	thermostatServer := &server{
		DB:          db,
		RobotClient: robotClient,
	}
	thermostatserver.RegisterThermostatServerServer(grpcServer, thermostatServer)
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
}
