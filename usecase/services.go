package usecase

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mrbenshef/SmartHomeAdapters/robotserver/robotserver"
	"google.golang.org/grpc"
)

const robotServerURL = "robotserver:8080"

// ConnectToRobotClient connects to the robot server and returns a client
// or an error if the connection could not be established.
func ConnectToRobotClient() (robotserver.RobotServerClient, error) {
	// connect to robotserver
	robotserverConn, err := grpc.Dial(robotServerURL, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer robotserverConn.Close()

	// return client
	client := robotserver.NewRobotServerClient(robotserverConn)
	return client, nil
}

// ConnectToDatabase connects to the postgres database for this usecase
// or an error if the connection could not be established.
func ConnectToDatabase() (*sql.DB, error) {
	var (
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		database = os.Getenv("DB_DATABASE")
		url      = os.Getenv("DB_URL")
	)

	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, url, database)
	log.Printf("Connecting to database with \"%s\"\n", connectionString)

	// connect to database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to postgres: %v", err)
	}

	// test connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %+v\n", db.Stats())
	return db, nil
}
