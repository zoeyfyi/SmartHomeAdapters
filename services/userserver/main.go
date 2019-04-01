//go:generate protoc --go_out=plugins=grpc:. ./userserver/userserver.proto
package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"regexp"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/userserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// bcryptRounds number of rounds bcrypt uses to hash passwords
const bcryptRounds = 10

type server struct {
	DB *sql.DB
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func (s *server) Register(ctx context.Context, request *userserver.RegisterRequest) (*empty.Empty, error) {
	// check fields are correct
	if request.Name == "" {
		return nil, status.New(codes.InvalidArgument, "Name is blank").Err()
	}
	if request.Email == "" {
		return nil, status.New(codes.InvalidArgument, "Email is blank").Err()
	}
	if !isEmailValid(request.Email) {
		return nil, status.Newf(codes.InvalidArgument, "Email \"%s\" is invalid", request.Email).Err()
	}
	if request.Password == "" {
		return nil, status.New(codes.InvalidArgument, "Password is blank").Err()
	}
	if len(request.Password) < 8 {
		return nil, status.New(codes.InvalidArgument, "Password is less than 8 characters").Err()
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcryptRounds)
	if err != nil {
		log.Printf("hash password failed: %v", err)
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	// insert user into database
	_, err = s.DB.Exec(
		"INSERT INTO users(username, email, password) VALUES($1, $2, $3)",
		request.Name,
		request.Email,
		hash,
	)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			// email already exists
			return nil, status.Newf(codes.AlreadyExists, "a user with email \"%s\" already exists", request.Email).Err()
		}

		log.Printf("failed to insert user into database: %v", err)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &empty.Empty{}, nil
}

func (s *server) CheckCredentials(ctx context.Context, credentials *userserver.Credentials) (*userserver.User, error) {
	// get email/hash from database
	var (
		id   string
		hash string
	)
	row := s.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", credentials.Email)
	err := row.Scan(&id, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Newf(codes.NotFound, "user with email \"%s\" does not exist", credentials.Email).Err()
		}

		log.Printf("error scanning database: %v", err)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(credentials.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, status.New(codes.InvalidArgument, "password incorrect").Err()
		}

		log.Printf("error comparing hashes: %v", err)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &userserver.User{
		Id: id,
	}, nil
}

func (s *server) GetUserID(ctx context.Context, email *userserver.Email) (*userserver.User, error) {
	// get email/hash from database
	var (
		id string
	)
	row := s.DB.QueryRow("SELECT id FROM users WHERE email = $1", email.Email)
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Newf(codes.NotFound, "user with email \"%s\" does not exist", email.Email).Err()
		}

		log.Printf("error scanning database: %v", err)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &userserver.User{
		Id: id,
	}, nil
}

func (s *server) GetUserByID(ctx context.Context, userID *userserver.UserId) (*userserver.User, error) {
	// get email/hash from database
	var (
		name string
	)
	row := s.DB.QueryRow("SELECT username FROM users WHERE id = $1", userID.Id)
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Newf(codes.NotFound, "user with id \"%s\" does not exist", userID.Id).Err()
		}

		log.Printf("error scanning database: %v", err)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &userserver.User{
		Id:   userID.Id,
		Name: name,
	}, nil
}

func main() {
	log.Println("Server starting")

	// connect to database
	db, err := microservice.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	defer db.Close()

	// test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %v+\n", db.Stats())

	// start grpc server
	grpcServer := grpc.NewServer()
	userServer := &server{DB: db}
	userserver.RegisterUserServerServer(grpcServer, userServer)
	lis, err := net.Listen("tcp", "userserver:80")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
