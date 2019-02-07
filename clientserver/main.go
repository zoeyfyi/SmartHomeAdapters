package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/julienschmidt/httprouter"
	"github.com/mrbenshef/SmartHomeAdapters/infoserver/infoserver"
	"google.golang.org/grpc"
)

var (
	infoserverClient infoserver.InfoServerClient
)

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

type Error struct {
	Error string `json:"error"`
}

type Robot struct {
	ID            string `json:"id"`
	Nickname      string `json:"nickname"`
	RobotType     string `json:"robotType"`
	InterfaceType string `json:"interfaceType"`
}

func robotsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	stream, err := infoserverClient.GetRobots(context.Background(), &empty.Empty{})
	if err != nil {
		json.NewEncoder(w).Encode(Error{err.Error()})
		return
	}

	var robots []Robot
	for {
		robot, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			json.NewEncoder(w).Encode(Error{err.Error()})
			return
		}
		robots = append(robots, Robot{
			ID:            robot.Id,
			Nickname:      robot.Nickname,
			RobotType:     robot.RobotType,
			InterfaceType: robot.InterfaceType,
		})
	}

	json.NewEncoder(w).Encode(robots)
}

func robotHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := fmt.Sprintf("http://infoserver/robot/%s", ps.ByName("id"))
	proxy(http.MethodGet, url, w, r)
}

func toggleHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := fmt.Sprintf("http://infoserver/robot/%s/toggle/%s", ps.ByName("id"), ps.ByName("value"))
	proxy(http.MethodPatch, url, w, r)
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
