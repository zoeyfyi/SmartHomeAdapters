package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/julienschmidt/httprouter"
	"github.com/mrbenshef/SmartHomeAdapters/infoserver/infoserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	infoserverClient infoserver.InfoServerClient
)

// HTTPStatusFromCode converts a gRPC error code into the corresponding HTTP response status.
func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	log.Printf("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}

type ErrorResponce struct {
	Error  string `json:"error"`
	Code   int    `json:"code"`
	Status int    `json:"status"`
}

// HTTPError transforms an error into a JSON responce
func HTTPError(w http.ResponseWriter, err error) {
	log.Printf("The following error occored: %v", err)

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	code := int(s.Code())
	status := HTTPStatusFromCode(s.Code())

	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(ErrorResponce{
		Error:  s.Message(),
		Code:   code,
		Status: status,
	})
	if err != nil {
		log.Printf("Failed to encode error responce: %v", err)
	}
}

type restError struct {
	Error string `json:"error"`
}

func httpWriteError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(restError{msg})
}

func pingHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("pong"))
}

// proxy forwards the request to a different url
func proxy(method string, url string, w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(method, url, r.Body)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.WriteHeader(resp.StatusCode)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	w.Write(buf.Bytes())
}

func registerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	proxy(http.MethodPost, "http://userserver/register", w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	proxy(http.MethodGet, "http://userserver/login", w, r)
}

type Robot struct {
	ID            string `json:"id"`
	Nickname      string `json:"nickname"`
	RobotType     string `json:"robotType"`
	InterfaceType string `json:"interfaceType"`
	Status        Status `json:"status,omitempty"`
}

type Status interface {
	status()
}

type ToggleStatus struct {
	Value bool `json:"value"`
}

func (s ToggleStatus) status() {}

type RangeStatus struct {
	Min     int `json:"min"`
	Max     int `json:"max"`
	Current int `json:"current"`
}

func (s RangeStatus) status() {}

func robotsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	stream, err := infoserverClient.GetRobots(context.Background(), &empty.Empty{})
	if err != nil {
		HTTPError(w, err)
		return
	}

	var robots []Robot
	for {
		robot, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			HTTPError(w, err)
			return
		}
		robots = append(robots, Robot{
			ID:            robot.Id,
			Nickname:      robot.Nickname,
			RobotType:     robot.RobotType,
			InterfaceType: robot.InterfaceType,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(robots)
}

func robotHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Printf("Getting robot with id %s", id)

	robot, err := infoserverClient.GetRobot(context.Background(), &infoserver.RobotQuery{Id: id})
	if err != nil {
		HTTPError(w, err)
		return
	}

	// convert robot status
	var status Status
	switch robotStatus := robot.RobotStatus.(type) {
	case *infoserver.Robot_ToggleStatus:
		status = ToggleStatus{
			Value: robotStatus.ToggleStatus.Value,
		}
	case *infoserver.Robot_RangeStatus:
		status = RangeStatus{
			Min:     int(robotStatus.RangeStatus.Min),
			Max:     int(robotStatus.RangeStatus.Max),
			Current: int(robotStatus.RangeStatus.Current),
		}
	default:
		panic(fmt.Sprintf("%T is not a valid robot status", robotStatus))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Robot{
		ID:            robot.Id,
		Nickname:      robot.Nickname,
		RobotType:     robot.RobotType,
		InterfaceType: robot.InterfaceType,
		Status:        status,
	})
}

func toggleHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	value := ps.ByName("value")

	log.Printf("Toggling robot %s", id)

	toggleValue, err := strconv.ParseBool(value)
	if err != nil {
		HTTPError(w, errors.New("Toggle value is not a boolean, should be either \"true\" or \"false\""))
		return
	}

	_, err = infoserverClient.ToggleRobot(context.Background(), &infoserver.ToggleRequest{
		Id:    id,
		Value: toggleValue,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}
}

func createRouter() *httprouter.Router {
	router := httprouter.New()

	// register routes
	router.GET("/ping", pingHandler)
	router.POST("/register", registerHandler)
	router.POST("/login", loginHandler)
	router.GET("/robots", robotsHandler)
	router.GET("/robot/:id", robotHandler)
	router.PATCH("/robot/:id/toggle/:value", toggleHandler)

	return router
}

func main() {
	log.Println("Server starting")

	// connect to infoserver
	infoserverConn, err := grpc.Dial("infoserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer infoserverConn.Close()
	infoserverClient = infoserver.NewInfoServerClient(infoserverConn)

	// start server
	if err := http.ListenAndServe(":80", createRouter()); err != nil {
		panic(err)
	}
}
