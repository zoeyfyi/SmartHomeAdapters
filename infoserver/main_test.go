package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	gock "gopkg.in/h2non/gock.v1"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/infoserver/infoserver"
	"google.golang.org/grpc"
)

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	// start test gRPC server
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	infoserver.RegisterInfoServerServer(s, &server{
		DB: getDb(),
	})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	os.Exit(m.Run())
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}
func TestGetRobots(t *testing.T) {
	expectedRobots := []*infoserver.Robot{
		&infoserver.Robot{
			Id:            "123abc",
			Nickname:      "testLightbot",
			RobotType:     "switch",
			InterfaceType: "toggle",
		},
		&infoserver.Robot{
			Id:            "T2D2",
			Nickname:      "testThermoBot",
			RobotType:     "thermostat",
			InterfaceType: "range",
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := infoserver.NewInfoServerClient(conn)
	stream, err := client.GetRobots(context.Background(), &empty.Empty{})
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
	cases := []struct {
		id            string
		expectedRobot *infoserver.Robot
		gockSetup     func()
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
						Value: false,
					},
				},
			},
			gockSetup: func() {
				gock.New("http://switchserver").
					Get("/123abc").
					Reply(200).
					JSON(map[string]interface{}{"isOn": false})
			},
		},
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
			gockSetup: func() {
				gock.New("http://switchserver").
					Get("/123abc").
					Reply(200).
					JSON(map[string]interface{}{"isOn": true})
			},
		},
	}

	for _, c := range cases {
		c.gockSetup()
		defer gock.Off()

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		client := infoserver.NewInfoServerClient(conn)
		robot, err := client.GetRobot(context.Background(), &infoserver.RobotQuery{Id: c.id})
		if err != nil {
			t.Fatalf("Could not get robot: %v", err)
		}

		if !reflect.DeepEqual(robot, c.expectedRobot) {
			t.Fatalf("Expected: %+v, Got: %+v", c.expectedRobot, robot)
		}
	}
}

func TestGetRobotWithInvalidID(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := infoserver.NewInfoServerClient(conn)
	robot, err := client.GetRobot(context.Background(), &infoserver.RobotQuery{Id: "invalidid"})

	status, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Expected grpc status error, got %v %T", err, err)
	}

	if status.Message() != "No robot with ID \"invalidid\"" {
		t.Errorf("Expected \"No robot with ID \"invalidid\"\" error message, got: %s", status.Message())
	}

	if robot != nil {
		t.Errorf("Robot was not nil")
	}
}
