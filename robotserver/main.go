//go:generate protoc --go_out=plugins=grpc:. ./robotserver/robotserver.proto
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/julienschmidt/httprouter"
	"github.com/mrbenshef/SmartHomeAdapters/robotserver/robotserver"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

type robotserverKey string

const dbKey robotserverKey = "db"

type server struct {
	DB *sql.DB
}

func (s *server) SetServo(ctx context.Context, request *robotserver.ServoRequest) (*empty.Empty, error) {
	log.Printf("setting servo to %d\n", request.Angle)

	if request.Angle > 180 {
		return nil, status.Newf(codes.InvalidArgument, "%d is to large, must be <= 180", request.Angle).Err()
	} else if request.Angle < 0 {
		return nil, status.Newf(codes.InvalidArgument, "%d is to small, must be >= 0", request.Angle).Err()
	}

	err := submitCommand(s.DB, int(request.Angle), int(request.Delay))
	if err != nil {
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	return &empty.Empty{}, nil
}

type command struct {
	id           string
	angle        int
	isCompleted  bool
	submitTime   time.Time
	completeTime *time.Time
	delayTime    int
}

func submitCommand(db *sql.DB, angle int, delay int) error {
	_, err := db.Exec("INSERT INTO command(angle, isCompleted, submitTime, completeTime, delayTime) VALUES($1, $2, $3, $4, $5)",
		angle,
		false,
		time.Now(),
		nil,
		delay,
	)

	return err
}

func getCommands(db *sql.DB) ([]command, error) {
	rows, err := db.Query("SELECT id, angle, isCompleted, submitTime, completeTime, delayTime FROM command WHERE isCompleted = FALSE ORDER BY submitTime ASC")
	if err != nil {
		log.Printf("Failed to query commands: %v", err)
		return nil, err
	}

	commands := []command{}

	for rows.Next() {
		var cmd command

		err := rows.Scan(&cmd.id, &cmd.angle, &cmd.isCompleted, &cmd.submitTime, &cmd.completeTime, &cmd.delayTime)
		if err != nil {
			log.Printf("Failed to scan command: %v", err)
			return nil, err
		}

		commands = append(commands, cmd)
	}

	return commands, nil
}

func acknowledgeCommand(db *sql.DB, cmdID string) error {
	res, err := db.Exec("UPDATE command SET isCompleted = TRUE WHERE id = $1", cmdID)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		return errors.New("Failed to acknowledge command")
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

func getCommandsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("Getting commands\n")

	db := r.Context().Value(dbKey).(*sql.DB)

	commands, err := getCommands(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	// build sequence of commands string
	var cmdString string
	for _, cmd := range commands {
		cmdString += fmt.Sprintf("(%s) servo %d;", cmd.id, cmd.angle)
		if cmd.delayTime > 0 {
			cmdString += fmt.Sprintf("(%s) delay %d", cmd.id, cmd.delayTime)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cmdString))
}

func acknowledgeCommandHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := r.Context().Value(dbKey).(*sql.DB)
	cmdID := ps.ByName("cmdID")

	err := acknowledgeCommand(db, cmdID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func connectionStr() string {
	var (
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		database = os.Getenv("DB_DATABASE")
		url      = os.Getenv("DB_URL")
	)

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

func dbProvider(h httprouter.Handle, db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r.WithContext(context.WithValue(nil, dbKey, db)), ps)
	}
}

func createRouter(db *sql.DB) *httprouter.Router {
	router := httprouter.New()
	router.GET("/:id/commands", dbProvider(getCommandsHandler, db))
	router.POST("/:id/acknowledge/:cmdID", dbProvider(acknowledgeCommandHandler, db))
	return router
}

func main() {
	db := getDb()
	defer db.Close()

	// test database
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	// start grpc server
	grpcServer := grpc.NewServer()
	robotServer := &server{DB: db}
	robotserver.RegisterRobotServerServer(grpcServer, robotServer)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting grpc server")

	go func() {
		grpcServer.Serve(lis)
	}()

	log.Println("Started grpc server, starting http server")

	// start REST server
	if err := http.ListenAndServe(":80", createRouter(db)); err != nil {
		panic(err)
	}
}
