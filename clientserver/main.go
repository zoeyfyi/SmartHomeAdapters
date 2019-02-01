package main

import (
	"encoding/json"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

type response struct {
	Id string `json:"id"`
	Nickname string `json:"nickname"`
	Interface struct {
		RobotType string `json:"type"`
		Min int `json:"min"`
		Max int `json:"max"`} `json:"interface"`}


func listRobotHandler(w http.ResponseWriter, r *http.Request){
	// Need to make a database connections and retrieve robots
	// TODO - Connect to database here


	resp := &response{
		Id : "123abc",
		Nickname : "myLightRobot",
		Interface: struct {
			RobotType string `json:"type"`
			Min       int `json:"min"`
			Max       int `json:"max"`
		} {RobotType: "lightSwitch", Min: 0, Max: 100} 	}

	jsonResponse, err := json.Marshal(resp)


	if (err != nil){
		w.Write([]byte("Something went wrong with the database. Please try again."))
		panic(err)
	}

	w.Write(jsonResponse)


}

func main() {
	// register routes
	// mux := http.NewServeMux()
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/robots", listRobotHandler)

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
