package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mrbenshef/SmartHomeAdapters/microservice/userserver"
	"github.com/ory/dockertest"
)

var testServer *server

func clearDatabase(t *testing.T) {
	_, err := testServer.DB.Exec("DELETE FROM users")
	if err != nil {
		t.Errorf("Error clearing database: %v", err)
	}
}

func TestMain(m *testing.M) {
	// connect to docker
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// start infodb
	resource, err := pool.Run("smarthomeadapters/userdb", "latest", []string{"POSTGRES_PASSWORD=password"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// wait till db is up
	if err = pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), dbDatabase))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dbURL = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))
	testServer = &server{DB: getDb()}

	exitCode := m.Run()

	pool.Purge(resource)

	os.Exit(exitCode)
}

func TestRegisterFieldValidation(t *testing.T) {
	clearDatabase(t)

	cases := []struct {
		request       *userserver.RegisterRequest
		expectedError string
	}{
		{
			&userserver.RegisterRequest{
				Email:    "",
				Password: "password",
			},
			"rpc error: code = InvalidArgument desc = Email is blank",
		},
		{
			&userserver.RegisterRequest{
				Email:    "foo",
				Password: "password",
			},
			"rpc error: code = InvalidArgument desc = Email \"foo\" is invalid",
		},
		{
			&userserver.RegisterRequest{
				Email:    "foo@bar.com",
				Password: "",
			},
			"rpc error: code = InvalidArgument desc = Password is blank",
		},
		{
			&userserver.RegisterRequest{
				Email:    "foo@bar.com",
				Password: "pass",
			},
			"rpc error: code = InvalidArgument desc = Password is less than 8 characters",
		},
	}

	for _, c := range cases {
		_, err := testServer.Register(context.Background(), c.request)

		if err == nil {
			t.Fatal("Expected non-nil error")
		}

		if err.Error() != c.expectedError {
			t.Errorf("Expected error: %s, got error: %+v", c.expectedError, err)
		}
	}
}

func TestSuccessfullRegistration(t *testing.T) {
	clearDatabase(t)

	user, err := testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	if err != nil {
		t.Errorf("Error registering user: %v", err)
	}

	if user == nil {
		t.Errorf("Expected user to be non-nil")
	}
}

func TestRegisterDuplicateEmails(t *testing.T) {
	clearDatabase(t)

	testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	user, err := testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	expectedError := "rpc error: code = AlreadyExists desc = A user with email \"foo@email.com\" already exists"

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, got error: %+v", expectedError, err)
	}

	if user != nil {
		t.Errorf("Expected user to be nil")
	}
}

func TestSuccessfullLogin(t *testing.T) {
	clearDatabase(t)

	testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	token, err := testServer.Login(context.Background(), &userserver.LoginRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	if err != nil {
		t.Errorf("Error login in user: %v", err)
	}

	if token == nil {
		t.Errorf("Expected user to be non-nil")
	}

	if token.Token == "" {
		t.Errorf("Expected non-empty token")
	}
}

func TestLoginFailure(t *testing.T) {
	clearDatabase(t)

	testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	cases := []struct {
		request       *userserver.LoginRequest
		expectedError string
	}{
		{
			&userserver.LoginRequest{
				Email:    "wrong email",
				Password: "password",
			},
			"rpc error: code = NotFound desc = User with email \"wrong email\" does not exist",
		},
		{
			&userserver.LoginRequest{
				Email:    "foo@bar.com",
				Password: "wrong password",
			},
			"rpc error: code = NotFound desc = User with email \"foo@bar.com\" does not exist",
		},
	}

	for _, c := range cases {
		user, err := testServer.Login(context.Background(), c.request)

		if err == nil {
			t.Fatal("Expected non-nil error")
		}

		if err.Error() != c.expectedError {
			t.Errorf("Expected error: %s, got error: %s", c.expectedError, err.Error())
		}

		if user != nil {
			t.Error("Expected nil user")
		}
	}
}

func TestSuccessfullAuthorization(t *testing.T) {
	clearDatabase(t)

	ctx := context.Background()

	testServer.Register(ctx, &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	token, _ := testServer.Login(ctx, &userserver.LoginRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	user, err := testServer.Authorize(ctx, &userserver.Token{
		Token: token.Token,
	})

	if err != nil {
		t.Errorf("Error authorizing user: %v", err)
	}

	if user == nil {
		t.Errorf("Expected user to be non-nil")
	}
}
