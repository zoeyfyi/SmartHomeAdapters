package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mrbenshef/SmartHomeAdapters/infoserver/infoserver"
	"github.com/mrbenshef/SmartHomeAdapters/userserver/userserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	infoserverClient infoserver.InfoServerClient
	userserverClient userserver.UserServerClient
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

func pingHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("pong"))
}

type registerBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var registerBody registerBody
	err := json.NewDecoder(r.Body).Decode(&registerBody)
	if err != nil {
		log.Printf("Invalid register JSON: %v", err)
		HTTPError(w, errors.New("Invalid JSON"))
		return
	}

	log.Printf("Registering user with email: %s", registerBody.Email)

	_, err = userserverClient.Register(context.Background(), &userserver.RegisterRequest{
		Email:    registerBody.Email,
		Password: registerBody.Password,
	})
	if err != nil {
		HTTPError(w, err)
	}
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponce struct {
	Token string `json:"token"`
}

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var loginBody loginBody
	err := json.NewDecoder(r.Body).Decode(&loginBody)
	if err != nil {
		log.Printf("Invalid register JSON: %v", err)
		HTTPError(w, errors.New("Invalid JSON"))
		return
	}

	log.Printf("Logging in user with email: %s", loginBody.Email)
	token, err := userserverClient.Login(context.Background(), &userserver.LoginRequest{
		Email:    loginBody.Email,
		Password: loginBody.Password,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&tokenResponce{
		Token: token.Token,
	})
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

	userId := r.Context().Value("userId").(string)

	stream, err := infoserverClient.GetRobots(context.Background(), &infoserver.RobotsQuery{UserId: userId})
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

	userId := r.Context().Value("userId").(string)

	robot, err := infoserverClient.GetRobot(context.Background(), &infoserver.RobotQuery{Id: id, UserId: userId})

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

type RegisterRobotBody struct {
	Nickname  string `json:"nickname"`
	RobotType string `json:"robotType"`
}

func registerRobotHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Printf("Registering robot with id %s", id)

	var registerRobotBody RegisterRobotBody

	err := json.NewDecoder(r.Body).Decode(&registerRobotBody)

	if err != nil {
		log.Printf("Invalid register JSON: %v", err)
		HTTPError(w, errors.New("Invalid JSON"))
		return
	}

	userId := r.Context().Value("userId").(string)

	registerQuery := infoserver.RegisterRobotQuery{Id: id, Nickname: registerRobotBody.Nickname, RobotType: registerRobotBody.RobotType, UserId: userId}

	_, err = infoserverClient.RegisterRobot(context.Background(), &registerQuery)

	if err != nil {
		HTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

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

type usecase struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Parameters []calibrationParameter `json:"parameters"`
}

type parameterDetails interface {
	parameterDetails()
}

type intParameter struct {
	Min     int `json:"min"`
	Max     int `json:"max"`
	Default int `json:"default"`
}

func (p intParameter) parameterDetails() {}

type booleanParameter struct {
	Default bool `json:"default"`
}

func (p booleanParameter) parameterDetails() {}

type calibrationParameter struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Details     parameterDetails `json:"details"`
}

// TODO: load from some kind of resource file
// we want this to extendable, so plugins can request a list of parameters
var usecases = map[string]usecase{
	"1": usecase{
		ID:   "1",
		Name: "Switch",
		Parameters: []calibrationParameter{
			calibrationParameter{
				Name:        "On Angle",
				Description: "Angle to turn the servo to turn the switch on",
				Type:        "int",
				Details: intParameter{
					Min:     90,
					Max:     180,
					Default: 100,
				},
			},
			calibrationParameter{
				Name:        "Off Angle",
				Description: "Angle to turn the servo to turn the switch off",
				Type:        "int",
				Details: intParameter{
					Min:     90,
					Max:     0,
					Default: 80,
				},
			},
		},
	},
	"2": usecase{
		ID:   "2",
		Name: "Thermostat",
		Parameters: []calibrationParameter{
			calibrationParameter{
				Name:        "Minimum tempreture",
				Description: "Minimum tempreture of the thermostat",
				Details: intParameter{
					Min:     200,
					Max:     400,
					Default: 293,
				},
			},
			calibrationParameter{
				Name:        "Maximum tempreture",
				Description: "Maximum tempreture of the thermostat",
				Details: intParameter{
					Min:     200,
					Max:     400,
					Default: 293,
				},
			},
			calibrationParameter{
				Name:        "Minimum tempreture angle",
				Description: "Angle of the minimum tempreture",
				Details: intParameter{
					Min:     0,
					Max:     180,
					Default: 80,
				},
			},
			calibrationParameter{
				Name:        "Maximum tempreture angle",
				Description: "Angle of the maximum tempreture",
				Details: intParameter{
					Min:     0,
					Max:     180,
					Default: 100,
				},
			},
		},
	},
}

func usecasesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// map -> list
	usecaseList := make([]usecase, 0, len(usecases))
	for _, usecase := range usecases {
		usecaseList = append(usecaseList, usecase)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usecaseList)
}

func usecaseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	usecase, ok := usecases[id]
	if !ok {
		err := status.Newf(codes.NotFound, "No usecase with id \"%s\"", id).Err()
		HTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usecase)
}

type parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func setCalibrationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	// decode json request
	var parameters []parameter
	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		log.Printf("Invalid calibration JSON: %v", err)
		HTTPError(w, errors.New("Invalid JSON"))
		return
	}

	// build infoserver calibration request
	request := &infoserver.CalibrationRequest{}
	request.Id = id
	for _, param := range parameters {
		request.Parameters = append(request.Parameters, &infoserver.CalibrationParameter{
			Name:  param.Name,
			Value: param.Value,
		})
	}

	// send request
	robot, err := infoserverClient.CalibrateRobot(context.Background(), request)
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

	// encode robot responce
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Robot{
		ID:            robot.Id,
		Nickname:      robot.Nickname,
		RobotType:     robot.RobotType,
		InterfaceType: robot.InterfaceType,
		Status:        status,
	})
}

func getCalibrationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	params, err := infoserverClient.GetCalibration(context.Background(), &infoserver.RobotQuery{
		Id: id,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}

	var parameters []parameter

	for _, p := range params.Parameters {
		parameters = append(parameters, parameter{
			Name:  p.Name,
			Value: p.Value,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(parameters)
}

func rangeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	value := ps.ByName("value")

	log.Printf("Setting range robot %s", id)

	rangeValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		HTTPError(w, errors.New("Toggle value is not a integer"))
		return
	}

	_, err = infoserverClient.RangeRobot(context.Background(), &infoserver.RangeRequest{
		Id:    id,
		Value: rangeValue,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}
}

func auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("token")
		user, err := userserverClient.Authorize(context.Background(), &userserver.Token{Token: token})

		if err != nil || user.Id == "" {
			log.Printf("Authentication error: %v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Printf("User ID is: %s", user.Id)
		h(w, r.WithContext(context.WithValue(nil, "userId", user.Id)), ps)
	}
}

func createRouter() *httprouter.Router {
	router := httprouter.New()

	// register routes
	router.GET("/ping", pingHandler)
	router.POST("/register", registerHandler)
	router.POST("/login", loginHandler)
	router.POST("/robot/:id", auth(registerRobotHandler))
	router.GET("/robots", auth(robotsHandler))
	router.GET("/robot/:id", auth(robotHandler))
	router.PATCH("/robot/:id/toggle/:value", auth(toggleHandler))
	router.PATCH("/robot/:id/range/:value", auth(rangeHandler))
	router.GET("/usecases", auth(usecasesHandler))
	router.GET("/usecase/:id", auth(usecaseHandler))
	router.GET("/robot/:id/calibration", auth(getCalibrationHandler))

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

	// connect to userserver
	userserverConn, err := grpc.Dial("userserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer userserverConn.Close()
	userserverClient = userserver.NewUserServerClient(userserverConn)

	// start server
	if err := http.ListenAndServe(":80", createRouter()); err != nil {
		panic(err)
	}
}
