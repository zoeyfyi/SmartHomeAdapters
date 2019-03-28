package main

import (
	"context"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/robotserver"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/thermostatserver"

	"google.golang.org/grpc/test/bufconn"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	empty "github.com/golang/protobuf/ptypes/empty"
)

var lis *bufconn.Listener
var servoRequests = []*robotserver.ServoRequest{}

func TestMain(m *testing.M) {
	username := os.Getenv("DB_USERNAME")
	if username != "temp" {
		log.Fatalf("Database username must be \"temp\", data will be wiped!")
	}

	db, err := microservice.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// start test gRPC server
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	thermostatserver.RegisterThermostatServerServer(s, &server{
		DB:          db,
		RobotClient: mockRobotServerClient{},
	})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	exitCode := m.Run()

	os.Exit(exitCode)
}

type mockRobotServerClient struct{}

func (c mockRobotServerClient) SetServo(
	ctx context.Context,
	in *robotserver.ServoRequest,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	servoRequests = append(servoRequests, in)
	return &empty.Empty{}, nil
}

func (c mockRobotServerClient) SetLED(
	ctx context.Context,
	in *robotserver.LEDRequest,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	panic("Call to SetLED!")
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

func TestSetThermostat(t *testing.T) {
	cases := []struct {
		thermostatRequest    *thermostatserver.SetThermostatRequest
		expectedServoRequest *robotserver.ServoRequest
	}{
		{
			thermostatRequest: &thermostatserver.SetThermostatRequest{
				Id:         "qwerty",
				Tempreture: 283,
				Unit:       "kelvin",
			},
			expectedServoRequest: &robotserver.ServoRequest{
				RobotId: "qwerty",
				Angle:   30,
			},
		},
		{
			thermostatRequest: &thermostatserver.SetThermostatRequest{
				Id:         "qwerty",
				Tempreture: 293,
				Unit:       "kelvin",
			},
			expectedServoRequest: &robotserver.ServoRequest{
				RobotId: "qwerty",
				Angle:   100,
			},
		},
		{
			thermostatRequest: &thermostatserver.SetThermostatRequest{
				Id:         "qwerty",
				Tempreture: 303,
				Unit:       "kelvin",
			},
			expectedServoRequest: &robotserver.ServoRequest{
				RobotId: "qwerty",
				Angle:   170,
			},
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := thermostatserver.NewThermostatServerClient(conn)

	for _, c := range cases {
		servoRequests = []*robotserver.ServoRequest{}

		stream, err := client.SetThermostat(ctx, c.thermostatRequest)

		for {
			status, err := stream.Recv()
			if err != nil {
				t.Errorf("Error receiving: %v", err)
			}

			if status.Status == thermostatserver.SetThermostatStatus_DONE {
				break
			}
		}

		if err != nil {
			t.Errorf("Expected nil error, got: %v", err)
		}

		if len(servoRequests) != 1 {
			t.Errorf("Expected single servo request, got: %+v (len = %d)", servoRequests, len(servoRequests))
		}

		if !reflect.DeepEqual(servoRequests[0], c.expectedServoRequest) {
			t.Errorf("Expected request: %+v, got %+v", c.expectedServoRequest, servoRequests[0])
		}
	}
}
