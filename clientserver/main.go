package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

func robotsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	proxy(http.MethodGet, "http://infoserver/robots", w, r)
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

	// start server
	if err := http.ListenAndServe(":80", createRouter()); err != nil {
		panic(err)
	}
}