//go:generate protoc --go_out=plugins=grpc:. ./userserver/userserver.proto
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/userserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// signingKey key for signing json web tokens
var signingKey = os.Getenv("JWT_SIGNING_KEY")

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
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	// insert user into database
	_, err = s.DB.Exec("INSERT INTO users(email, password) VALUES($1, $2)", request.Email, hash)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			// email already exists
			return nil, status.Newf(codes.AlreadyExists, "A user with email \"%s\" already exists", request.Email).Err()
		}

		log.Printf("Failed to insert user into database: %v", err)
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	return &empty.Empty{}, nil
}

func (s *server) Login(ctx context.Context, request *userserver.LoginRequest) (*userserver.Token, error) {
	// get email/hash from database
	var id int
	var hash string
	row := s.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", request.Email)
	err := row.Scan(&id, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Newf(codes.NotFound, "User with email \"%s\" does not exist", request.Email).Err()
		}

		log.Printf("Error scanning database: %v", err)
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(request.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, status.New(codes.InvalidArgument, "Password incorrect").Err()
		}

		log.Printf("Error comparing hashes: %v", err)
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	// create token
	expire := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expire,
		"id":  strconv.Itoa(id),
	})
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return nil, status.New(codes.Internal, "Internal error").Err()
	}

	return &userserver.Token{
		Token: tokenString,
	}, nil
}

func (s *server) Authorize(ctx context.Context, token *userserver.Token) (*userserver.User, error) {

	// parse token
	jwtToken, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		log.Printf("could not parse token: %v", err)
		return nil, status.New(codes.InvalidArgument, "Invalid token").Err()
	}

	// get claims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		log.Printf("Could not get token claims")
		return nil, status.New(codes.InvalidArgument, "Invalid token").Err()
	}

	// get ID from claims
	id, ok := claims["id"].(string)
	if !ok {
		log.Printf("Token did not have ID claim, claims: %+v", claims)
		return nil, status.New(codes.InvalidArgument, "Invalid token").Err()
	}

	return &userserver.User{
		Id: id,
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
	lis, err := net.Listen("tcp", "127.0.0.1:80")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
