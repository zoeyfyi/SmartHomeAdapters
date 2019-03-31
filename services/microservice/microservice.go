package microservice

//go:generate protoc --go_out=plugins=grpc:. ./infoserver/infoserver.proto
//go:generate protoc --go_out=plugins=grpc:. ./robotserver/robotserver.proto
//go:generate protoc --go_out=plugins=grpc:. ./switchserver/switchserver.proto
//go:generate protoc --go_out=plugins=grpc:. ./thermostatserver/thermostatserver.proto
//go:generate protoc --go_out=plugins=grpc:. ./userserver/userserver.proto
//go:generate protoc --go_out=plugins=grpc:. ./usecase-service/usecase-service.proto

import (
	"database/sql"
	"fmt"
	"os"
)

// ConnectToDB connects to a database using the credentials from enviroment variables
func ConnectToDB() (*sql.DB, error) {
	var (
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		database = os.Getenv("DB_DATABASE")
		url      = os.Getenv("DB_URL")
	)

	connectionStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, url, database)
	return sql.Open("postgres", connectionStr)
}
