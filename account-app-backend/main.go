package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: check email/password are valid
	
}

func main() {
	router := httprouter.New()

	router.POST("/login", loginHandler)

	// start server
	if err := http.ListenAndServe(":80", router); err != nil {
		panic(err)
	}
}
