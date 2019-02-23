package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ory/dockertest"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/mrbenshef/SmartHomeAdapters/infoserver/infoserver"
	"github.com/mrbenshef/SmartHomeAdapters/switchserver/switchserver"
	"google.golang.org/grpc"
)

var lis *bufconn.Listener

type mockSwitchClient struct{}

func (c mockSwitchClient) GetSwitch(ctx context.Context, query *switchserver.SwitchQuery, _ ...grpc.CallOption) (*switchserver.Switch, error) {
	if query.Id == "123abc" {
		return &switchserver.Switch{
			Id:   "123abc",
			IsOn: true,
		}, nil
	} else {
		return nil, status.New(codes.NotFound, "Switch does not exist").Err()
	}
}

func (c mockSwitchClient) AddSwitch(_ context.Context, _ *switchserver.AddSwitchRequest, _ ...grpc.CallOption) (*switchserver.Switch, error) {
	return nil, nil
}

func (c mockSwitchClient) RemoveSwitch(_ context.Context, _ *switchserver.RemoveSwitchRequest, _ ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (c mockSwitchClient) SetSwitch(_ context.Context, _ *switchserver.SetSwitchRequest, _ ...grpc.CallOption) (switchserver.SwitchServer_SetSwitchClient, error) {
	return nil, nil
}

func TestMain(m *testing.M) {
	// connect to docker
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// start infodb
	resource, err := pool.Run("smarthomeadapters/infodb", "latest", []string{"POSTGRES_PASSWORD=password"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	url = fmt.Sprintf("localhost:%s", resource.GetPort("5432/tcp"))

	// wait till db is up
	if err = pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), database))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// start test gRPC server
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	infoserver.RegisterInfoServerServer(s, &server{
		DB:           getDb(),
		SwitchClient: mockSwitchClient{},
	})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	exitCode := m.Run()

	pool.Purge(resource)

	os.Exit(exitCode)
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
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := infoserver.NewInfoServerClient(conn)
	stream, err := client.GetRobots(context.Background(), &infoserver.RobotsQuery{UserId:"1"})
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
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		client := infoserver.NewInfoServerClient(conn)
		robot, err := client.GetRobot(context.Background(), &infoserver.RobotQuery{Id: c.id, UserId:"1"})

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
	robot, err := client.GetRobot(context.Background(), &infoserver.RobotQuery{Id: "invalidid", UserId:"1"})

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
