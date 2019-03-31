package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/userserver"
)

var testServer *server

func clearDatabase(t *testing.T) {
	_, err := testServer.DB.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("Error clearing database: %v", err)
	}
}

func TestMain(m *testing.M) {
	username := os.Getenv("DB_USERNAME")
	if username != "temp" {
		log.Fatalf("Database username must be \"temp\", data will be wiped!")
	}

	db, err := microservice.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	testServer = &server{DB: db}

	exitCode := m.Run()

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

	_, err := testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}

	user, err := testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})

	expectedError := "rpc error: code = AlreadyExists desc = a user with email \"foo@email.com\" already exists"

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

func TestSuccessfullCheckCredentials(t *testing.T) {
	clearDatabase(t)

	_, err := testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}

	user, err := testServer.CheckCredentials(context.Background(), &userserver.Credentials{
		Email:    "foo@email.com",
		Password: "password",
	})

	if err != nil {
		t.Errorf("Error login in user: %v", err)
	}

	if user == nil {
		t.Errorf("Expected user to be non-nil")
	}
}

func TestCheckCredentialsFailure(t *testing.T) {
	clearDatabase(t)

	_, err := testServer.Register(context.Background(), &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}

	cases := []struct {
		request       *userserver.Credentials
		expectedError string
	}{
		{
			&userserver.Credentials{
				Email:    "wrong email",
				Password: "password",
			},
			"rpc error: code = NotFound desc = user with email \"wrong email\" does not exist",
		},
		{
			&userserver.Credentials{
				Email:    "foo@bar.com",
				Password: "wrong password",
			},
			"rpc error: code = NotFound desc = user with email \"foo@bar.com\" does not exist",
		},
	}

	for _, c := range cases {
		user, err := testServer.CheckCredentials(context.Background(), c.request)

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

func TestSuccessfullGetUserID(t *testing.T) {
	clearDatabase(t)

	ctx := context.Background()

	_, err := testServer.Register(ctx, &userserver.RegisterRequest{
		Email:    "foo@email.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}

	user, err := testServer.GetUserID(ctx, &userserver.Email{
		Email: "foo@email.com",
	})

	if err != nil {
		t.Errorf("Error authorizing user: %v", err)
	}

	if user == nil {
		t.Errorf("Expected user to be non-nil")
	}
}
