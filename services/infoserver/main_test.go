package main

import (
	"context"
	"database/sql"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/microservice"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/infoserver"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/switchserver"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/thermostatserver"
	"google.golang.org/grpc"
)

var lis *bufconn.Listener
var db *sql.DB

func resetDatabase(t *testing.T) {
	_, err := db.Exec("DROP TABLE robots")
	if err != nil {
		t.Fatalf("Error dropping table: %v", err)
	}

	dbSQL, err := ioutil.ReadFile("../infodb/init.sql")
	if err != nil {
		t.Fatalf("Error reading SQL: %v", err)
	}

	_, err = db.Exec(string(dbSQL))
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}
}

type mockSwitchClient struct{}

func (c mockSwitchClient) GetSwitch(
	ctx context.Context,
	query *switchserver.SwitchQuery,
	_ ...grpc.CallOption,
) (*switchserver.Switch, error) {
	if query.Id == "123abc" {
		return &switchserver.Switch{
			Id:   "123abc",
			IsOn: true,
		}, nil
	}

	return nil, status.New(codes.NotFound, "Switch does not exist").Err()
}

func (c mockSwitchClient) AddSwitch(
	_ context.Context,
	_ *switchserver.AddSwitchRequest,
	_ ...grpc.CallOption,
) (*switchserver.Switch, error) {
	return nil, nil
}

func (c mockSwitchClient) RemoveSwitch(
	_ context.Context,
	_ *switchserver.RemoveSwitchRequest,
	_ ...grpc.CallOption,
) (*empty.Empty, error) {
	return nil, nil
}

func (c mockSwitchClient) SetSwitch(
	_ context.Context,
	_ *switchserver.SetSwitchRequest,
	_ ...grpc.CallOption,
) (switchserver.SwitchServer_SetSwitchClient, error) {
	return nil, nil
}

func (c mockSwitchClient) CalibrateSwitch(
	_ context.Context,
	_ *switchserver.SwitchCalibrationParameters,
	_ ...grpc.CallOption,
) (*empty.Empty, error) {
	return nil, nil
}

type mockThermostatClient struct{}

func (c mockThermostatClient) GetThermostat(
	ctx context.Context,
	query *thermostatserver.ThermostatQuery,
	_ ...grpc.CallOption,
) (*thermostatserver.Thermostat, error) {
	if query.Id == "qwerty" {
		return &thermostatserver.Thermostat{
			Id:            "qwerty",
			MaxTempreture: 40,
			MinTempreture: 20,
			Tempreture:    30,
			MinAngle:      0,
			MaxAngle:      180,
		}, nil
	}

	return nil, status.New(codes.NotFound, "Thermostat does not exist").Err()
}

func (c mockThermostatClient) SetThermostat(
	_ context.Context,
	_ *thermostatserver.SetThermostatRequest,
	_ ...grpc.CallOption,
) (thermostatserver.ThermostatServer_SetThermostatClient, error) {
	return nil, nil
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

	infoserver.RegisterInfoServerServer(s, &server{
		DB:               db,
		SwitchClient:     mockSwitchClient{},
		ThermostatClient: mockThermostatClient{},
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

func TestGetRobots(t *testing.T) {
	resetDatabase(t)

	expectedRobots := []*infoserver.Robot{
		{
			Id:            "123abc",
			Nickname:      "testLightbot",
			RobotType:     "switch",
			InterfaceType: "toggle",
			RobotStatus: &infoserver.Robot_ToggleStatus{
				ToggleStatus: &infoserver.ToggleStatus{
					Value: true,
				},
			},
		},
		{
			Id:            "qwerty",
			Nickname:      "testThermoBot",
			RobotType:     "thermostat",
			InterfaceType: "range",
			RobotStatus: &infoserver.Robot_RangeStatus{
				RangeStatus: &infoserver.RangeStatus{
					Min:     20,
					Max:     40,
					Current: 30,
				},
			},
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := infoserver.NewInfoServerClient(conn)
	stream, err := client.GetRobots(context.Background(), &infoserver.RobotsQuery{UserId: "1"})
	if err != nil {
		t.Fatalf("Could not get robots: %v", err)
	}

	var robots []*infoserver.Robot
	for {
		robot, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to receive robot: %v", err)
		}
		robots = append(robots, robot)
	}

	if !reflect.DeepEqual(robots, expectedRobots) {
		t.Fatalf("Expected: %+v, Got: %+v", expectedRobots, robots)
	}

}

func TestGetRobotWithValidID(t *testing.T) {
	resetDatabase(t)

	cases := []struct {
		id            string
		expectedRobot *infoserver.Robot
	}{
		{
			id: "123abc",
			expectedRobot: &infoserver.Robot{
				Id:            "123abc",
				Nickname:      "testLightbot",
				RobotType:     "switch",
				InterfaceType: "toggle",
				RobotStatus: &infoserver.Robot_ToggleStatus{
					ToggleStatus: &infoserver.ToggleStatus{
						Value: true,
					},
				},
			},
		},
	}

	for _, c := range cases {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		client := infoserver.NewInfoServerClient(conn)
		robot, err := client.GetRobot(context.Background(), &infoserver.RobotQuery{Id: c.id, UserId: "1"})

		if err != nil {
			t.Fatalf("Could not get robot: %v", err)
		}

		if !reflect.DeepEqual(robot, c.expectedRobot) {
			t.Fatalf("Expected: %+v, Got: %+v", c.expectedRobot, robot)
		}
	}
}

func TestGetRobotWithInvalidID(t *testing.T) {
	resetDatabase(t)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := infoserver.NewInfoServerClient(conn)
	robot, err := client.GetRobot(context.Background(), &infoserver.RobotQuery{Id: "invalidid", UserId: "1"})

	status, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Expected grpc status error, got %v %T", err, err)
	}

	expectedMessage := "Robot \"invalidid\" does not exist"

	if status.Message() != expectedMessage {
		t.Errorf("Expected %s error message, got: %s", expectedMessage, status.Message())
	}

	if robot != nil {
		t.Errorf("Robot was not nil")
	}
}

func TestSetRobotUsecase(t *testing.T) {
	resetDatabase(t)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := infoserver.NewInfoServerClient(conn)
	_, err = client.SetUsecase(ctx, &infoserver.SetUsecaseRequest{
		Id:      "qwerty",
		UserId:  "1",
		Usecase: "switch",
	})

	expectedError := "rpc error: code = NotFound desc = Switch does not exist"
	if err.Error() != expectedError {
		t.Fatalf("Expected error: %s, got error: %v", expectedError, err)
	}
}
