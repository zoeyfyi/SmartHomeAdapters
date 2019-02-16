package usecaseapi

import (
	"github.com/mrbenshef/SmartHomeAdapters/robotserver/robotserver"
	"google.golang.org/grpc"
)

const robotServerURL = "robotserver:8080"

// ConnectToRobotClient connects to the robot server and returns a client
// or an error if the connection could not be established.
func ConnectToRobotClient() (robotserver.RobotServerClient, error) {
	// connect to robotserver
	robotserverConn, err := grpc.Dial(robotServerURL, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer robotserverConn.Close()

	// return client
	client := robotserver.NewRobotServerClient(robotserverConn)
	return client, nil
}
