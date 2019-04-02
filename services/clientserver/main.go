package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/julienschmidt/httprouter"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/infoserver"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/userserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type clientServerKey string

var userIDKey clientServerKey = "userID"

var (
	infoserverClient infoserver.InfoServerClient
	userserverClient userserver.UserServerClient
)

var (
	errInvalidJSON  = status.New(codes.InvalidArgument, "invalid JSON").Err()
	errFailedEncode = status.New(codes.Internal, "failed to encode responce").Err()
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

type errorResponce struct {
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
	err = json.NewEncoder(w).Encode(errorResponce{
		Error:  s.Message(),
		Code:   code,
		Status: status,
	})
	if err != nil {
		log.Printf("Failed to encode error responce: %v", err)
		HTTPError(w, errFailedEncode)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Printf("Error writing responce: %v", err)
	}
}

type registerBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body registerBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Invalid register JSON: %v", err)
		HTTPError(w, errInvalidJSON)
		return
	}

	log.Printf("Registering user with email: %s", body.Email)

	_, err = userserverClient.Register(context.Background(), &userserver.RegisterRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		HTTPError(w, err)
	}
}

type userResponce struct {
	Name string `json:"name"`
}

func userHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Context().Value(userIDKey).(string)
	log.Printf("getting user with id %s", userID)

	user, err := userserverClient.GetUserByID(context.Background(), &userserver.UserId{
		Id: userID,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(userResponce{
		Name: user.Name,
	})
	if err != nil {
		log.Printf("failed to encode error responce: %v", err)
		HTTPError(w, errFailedEncode)
	}
}

type robot struct {
	ID            string      `json:"id"`
	Nickname      string      `json:"nickname"`
	RobotType     string      `json:"robotType"`
	InterfaceType string      `json:"interfaceType"`
	Status        robotStatus `json:"status,omitempty"`
}

type robotStatus interface {
	status()
}

type toggleStatus struct {
	Value bool `json:"value"`
}

func (s toggleStatus) status() {}

type rangeStatus struct {
	Min     int `json:"min"`
	Max     int `json:"max"`
	Current int `json:"current"`
}

func (s rangeStatus) status() {}

func robotsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value(userIDKey).(string)

	stream, err := infoserverClient.GetRobots(context.Background(), &infoserver.RobotsQuery{
		UserId: userID,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}

	var robots []robot
	for {
		rbt, recvErr := stream.Recv()
		if recvErr == io.EOF {
			break
		}
		if recvErr != nil {
			HTTPError(w, recvErr)
			return
		}

		// convert robot status
		var status robotStatus
		switch robotStatus := rbt.RobotStatus.(type) {
		case *infoserver.Robot_ToggleStatus:
			status = toggleStatus{
				Value: robotStatus.ToggleStatus.Value,
			}
		case *infoserver.Robot_RangeStatus:
			status = rangeStatus{
				Min:     int(robotStatus.RangeStatus.Min),
				Max:     int(robotStatus.RangeStatus.Max),
				Current: int(robotStatus.RangeStatus.Current),
			}
		default:
			panic(fmt.Sprintf("%T is not a valid robot status", robotStatus))
		}

		robots = append(robots, robot{
			ID:            rbt.Id,
			Nickname:      rbt.Nickname,
			RobotType:     rbt.RobotType,
			InterfaceType: rbt.InterfaceType,
			Status:        status,
		})
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(robots)
	if err != nil {
		log.Printf("Failed to encode error responce: %v", err)
		HTTPError(w, errFailedEncode)
	}
}

func robotHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Printf("Getting robot with id %s", id)

	userID := r.Context().Value(userIDKey).(string)

	rbt, err := infoserverClient.GetRobot(context.Background(), &infoserver.RobotQuery{Id: id, UserId: userID})

	if err != nil {
		HTTPError(w, err)
		return
	}

	// convert robot status
	var status robotStatus
	switch robotStatus := rbt.RobotStatus.(type) {
	case *infoserver.Robot_ToggleStatus:
		status = toggleStatus{
			Value: robotStatus.ToggleStatus.Value,
		}
	case *infoserver.Robot_RangeStatus:
		status = rangeStatus{
			Min:     int(robotStatus.RangeStatus.Min),
			Max:     int(robotStatus.RangeStatus.Max),
			Current: int(robotStatus.RangeStatus.Current),
		}
	default:
		panic(fmt.Sprintf("%T is not a valid robot status", robotStatus))
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(robot{
		ID:            rbt.Id,
		Nickname:      rbt.Nickname,
		RobotType:     rbt.RobotType,
		InterfaceType: rbt.InterfaceType,
		Status:        status,
	})
	if err != nil {
		log.Printf("Failed to encode error responce: %v", err)
		HTTPError(w, errFailedEncode)
	}
}

type registerRobotBody struct {
	Nickname  string `json:"nickname"`
	RobotType string `json:"robotType"`
}

func registerRobotHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	log.Printf("Registering robot with id %s", id)

	var body registerRobotBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Printf("Invalid register JSON: %v", err)
		HTTPError(w, errInvalidJSON)
		return
	}

	userID := r.Context().Value(userIDKey).(string)

	registerQuery := infoserver.RegisterRobotQuery{
		Id:        id,
		Nickname:  body.Nickname,
		RobotType: body.RobotType,
		UserId:    userID,
	}

	_, err = infoserverClient.RegisterRobot(context.Background(), &registerQuery)
	if err != nil {
		HTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func unregisterRobotHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	userID := r.Context().Value(userIDKey).(string)
	log.Printf("unregistering robot with id %s", id)

	_, err := infoserverClient.UnregisterRobot(context.Background(), &infoserver.UnregisterRobotQuery{
		Id:     id,
		UserId: userID,
	})
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
		HTTPError(w, errors.New("toggle value is not a boolean, should be either \"true\" or \"false\""))
		return
	}

	_, err = infoserverClient.ToggleRobot(context.Background(), &infoserver.ToggleRequest{
		Id:     id,
		UserId: r.Context().Value(userIDKey).(string),
		Value:  toggleValue,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}
}

type usecase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func usecasesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	stream, err := infoserverClient.GetUsecases(context.Background(), &empty.Empty{})
	if err != nil {
		HTTPError(w, err)
		return
	}

	usecases := make([]usecase, 0)

	for {
		uc, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			HTTPError(w, err)
			return
		}

		usecases = append(usecases, usecase{
			Name:        uc.Name,
			Description: uc.Description,
		})
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(usecases)
	if err != nil {
		log.Printf("Failed to encode error responce: %v", err)
		HTTPError(w, errFailedEncode)
	}
}

type setParameterRequest struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func setCalibrationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	robotID := ps.ByName("id")

	// decode json request
	var setParameter setParameterRequest
	err := json.NewDecoder(r.Body).Decode(&setParameter)
	if err != nil {
		log.Printf("Invalid calibration JSON: %v", err)
		HTTPError(w, errInvalidJSON)
		return
	}

	// build infoserver calibration request
	request := &infoserver.CalibrationRequest{
		Id:      setParameter.ID,
		RobotId: robotID,
		UserId:  r.Context().Value(userIDKey).(string),
	}
	log.Printf("calibration request: %+v", request)

	switch setParameter.Type {
		case "bool":
		value, err := strconv.ParseBool(setParameter.Value)
			if err != nil {
			HTTPError(w, fmt.Errorf("value should be either \"true\" or \"false\""))
				return
			}
		log.Printf("value: %t", value)
		request.Value = &infoserver.CalibrationRequest_BoolValue{
				BoolValue: value,
			}
		case "int":
		value, err := strconv.ParseInt(setParameter.Value, 10, 64)
			if err != nil {
			HTTPError(w, fmt.Errorf("value must be an integer"))
				return
			}
		log.Printf("value: %d", value)
		request.Value = &infoserver.CalibrationRequest_IntValue{
				IntValue: value,
			}
		default:
		HTTPError(w, fmt.Errorf("\"%s\" is not a recognized parameter type", setParameter.Type))
			return
		}

	// send request
	rbt, err := infoserverClient.CalibrateRobot(context.Background(), request)
	if err != nil {
		HTTPError(w, err)
		return
	}

	// convert robot status
	var status robotStatus
	switch robotStatus := rbt.RobotStatus.(type) {
	case *infoserver.Robot_ToggleStatus:
		status = toggleStatus{
			Value: robotStatus.ToggleStatus.Value,
		}
	case *infoserver.Robot_RangeStatus:
		status = rangeStatus{
			Min:     int(robotStatus.RangeStatus.Min),
			Max:     int(robotStatus.RangeStatus.Max),
			Current: int(robotStatus.RangeStatus.Current),
		}
	default:
		panic(fmt.Sprintf("%T is not a valid robot status", robotStatus))
	}

	// encode robot responce
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(robot{
		ID:            rbt.Id,
		Nickname:      rbt.Nickname,
		RobotType:     rbt.RobotType,
		InterfaceType: rbt.InterfaceType,
		Status:        status,
	})
	if err != nil {
		log.Printf("Failed to encode error responce: %v", err)
		HTTPError(w, errFailedEncode)
	}
}

type parameter struct {
	ID          string  `json:"id"`
	Name        *string `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Details     details `json:"details"`
}

type details interface {
	details()
}

type intDetails struct {
	Min     int64 `json:"min"`
	Max     int64 `json:"max"`
	Default int64 `json:"default"`
	Current int64 `json:"current"`
}

func (d intDetails) details() {}

type boolDetails struct {
	Default bool `json:"default"`
	Current bool `json:"current"`
}

func (d boolDetails) details() {}

func getCalibrationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	userID := r.Context().Value(userIDKey).(string)

	params, err := infoserverClient.GetCalibration(context.Background(), &infoserver.RobotQuery{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}

	var parameters []parameter

	for _, p := range params.Parameters {
		param := parameter{
			ID:          p.Id,
			Name:        &p.Name,
			Description: p.Description,
			Type:        p.Type,
		}

		switch details := p.Details.(type) {
		case *infoserver.CalibrationParameter_BoolDetails:
			param.Details = boolDetails{
				Default: details.BoolDetails.Default,
				Current: details.BoolDetails.Current,
			}
		case *infoserver.CalibrationParameter_IntDetails:
			param.Details = intDetails{
				Min:     details.IntDetails.Min,
				Max:     details.IntDetails.Max,
				Default: details.IntDetails.Default,
				Current: details.IntDetails.Current,
			}
		}

		parameters = append(parameters, param)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(parameters)
	if err != nil {
		log.Printf("Failed to encode error responce: %v", err)
		HTTPError(w, errFailedEncode)
	}
}

func rangeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	value := ps.ByName("value")

	log.Printf("Setting range robot %s", id)

	rangeValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		HTTPError(w, errors.New("toggle value is not a integer"))
		return
	}

	_, err = infoserverClient.RangeRobot(context.Background(), &infoserver.RangeRequest{
		Id:     id,
		UserId: r.Context().Value(userIDKey).(string),
		Value:  rangeValue,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}
}

type renameRequest struct {
	Nickname string `json:"nickname"`
}

func renameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	log.Printf("renaming robot \"%s\"", id)

	// decode json request
	var request renameRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Invalid rename request JSON: %v", err)
		HTTPError(w, errInvalidJSON)
		return
	}

	log.Printf("renaming robot \"%s\" to \"%s\"", id, request.Nickname)

	_, err = infoserverClient.RenameRobot(context.Background(), &infoserver.RenameRobotRequest{
		Id:          id,
		UserId:      r.Context().Value(userIDKey).(string),
		NewNickname: request.Nickname,
	})
	if err != nil {
		HTTPError(w, err)
		return
	}
}

type hydraIntrospect struct {
	Subject   string `json:"sub"`
	Active    bool   `json:"active"`
	Scope     string `json:"scope"`
	ClientID  string `json:"client_id"`
	Exp       int    `json:"exp"`
	Iat       int    `json:"iat"`
	Iss       string `json:"iss"`
	TokenType string `json:"token_type"`
}

func auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("token")
		// no
		// user, err := userserverClient.Authorize(context.Background(), &userserver.Token{Token: token})

		// send this to introspect

		urlThing := "https://hydra.halspals.co.uk/oauth2/introspect"

		// need to put json data here

		// need to fix redirect as well
		formdata := url.Values{
			"token": {token},
		}
		log.Printf("Sending token: %s to url %s", formdata, urlThing)

		resp, err := http.PostForm(urlThing, formdata)
		// parse the subject from the request and pass it on

		if err != nil {
			log.Printf("HTTP do error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var introspect hydraIntrospect

		err = json.NewDecoder(resp.Body).Decode(&introspect)

		if err != nil {
			log.Printf("JSON decode error %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if introspect.Subject == "" {
			log.Printf("Authentication error: %v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Printf("User ID is: %s", introspect.Subject)

		h(w, r.WithContext(context.WithValue(context.TODO(), userIDKey, introspect.Subject)), ps)
	}
}

func createRouter() *httprouter.Router {
	router := httprouter.New()

	// register routes
	router.GET("/ping", pingHandler)
	router.POST("/register", registerHandler)
	router.GET("/user", auth(userHandler))
	router.GET("/robots", auth(robotsHandler))
	router.GET("/robot/:id", auth(robotHandler))
	router.POST("/robot/:id", auth(registerRobotHandler))
	router.DELETE("/robot/:id", auth(unregisterRobotHandler))
	router.PATCH("/robot/:id/toggle/:value", auth(toggleHandler))
	router.PATCH("/robot/:id/range/:value", auth(rangeHandler))
	router.PATCH("/robot/:id/nickname", auth(renameHandler))
	router.GET("/usecases", usecasesHandler)
	router.GET("/robot/:id/calibration", auth(getCalibrationHandler))
	router.PUT("/robot/:id/calibration", auth(setCalibrationHandler))

	return router
}

func main() {
	log.Println("Server starting ")

	// TODO: FIX THIS
	// nolint
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// connect to infoserver
	infoserverConn, err := grpc.Dial("infoserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer func() {
		closeErr := infoserverConn.Close()
		if closeErr != nil {
			log.Fatalf("failed to close infoserver connection: %v", closeErr)
		}
	}()

	infoserverClient = infoserver.NewInfoServerClient(infoserverConn)

	// connect to userserver
	userserverConn, err := grpc.Dial("userserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer func() {
		closeErr := userserverConn.Close()
		if closeErr != nil {
			log.Fatalf("failed to close userserver connection: %v", closeErr)
		}
	}()

	userserverClient = userserver.NewUserServerClient(userserverConn)

	// start server
	if err := http.ListenAndServe(":80", createRouter()); err != nil {
		panic(err)
	}
}
