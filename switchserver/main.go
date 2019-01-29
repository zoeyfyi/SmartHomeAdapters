package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// database connection infomation
var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)

func connectionStr() string {
	if username == "" {
		username = "postgres"
	}
	if password == "" {
		password = "password"
	}
	if url == "" {
		url = "localhost:5432"
	}
	if database == "" {
		database = "postgres"
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, url, database)
}

func getDb() *sql.DB {
	log.Printf("Connecting to database with \"%s\"\n", connectionStr())
	db, err := sql.Open("postgres", connectionStr())
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	return db
}

type switchRobot struct {
	IsOn bool `json:"isOn"`
}

type restError struct {
	Error string `json:"error"`
}

func httpWriteError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(restError{msg})
}

// error messages for register route
const (
	ErrorInvalidJSON     = "Invalid JSON"
	ErrorIsOnMissing     = "Field \"isOn\" missing"
	ErrorRobotRegistered = "Robot is already registered"
	ErrorInvalidRobotID  = "Invalid robot ID"
)

type addSwitchBody struct {
	IsOn *bool `json:"isOn"`
}

// addSwitchHandler registers a robot as a switch robot
func addSwitchHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// decode body
		var addSwitchBody addSwitchBody
		err := json.NewDecoder(r.Body).Decode(&addSwitchBody)
		if err != nil {
			httpWriteError(w, ErrorInvalidJSON, http.StatusBadRequest)
			return
		}

		// check fields
		if addSwitchBody.IsOn == nil {
			httpWriteError(w, ErrorIsOnMissing, http.StatusBadRequest)
			return
		}

		// parse robot ID
		robotID, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			httpWriteError(w, ErrorInvalidRobotID, http.StatusBadRequest)
			return
		}

		// insert into database
		_, err = db.Exec("INSERT INTO switches(robotId, isOn) VALUES($1, $2)", robotID, *addSwitchBody.IsOn)
		if err != nil {
			if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
				// robot is already registered
				httpWriteError(w, ErrorRobotRegistered, http.StatusBadRequest)
			} else {
				log.Printf("Failed to insert user into database: %v", err)
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		// return success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(switchRobot{
			IsOn: *addSwitchBody.IsOn,
		})
	}
}

// error messages for remove switch routes
const (
	ErrorRobotNotRegistered = "Robot is not registered"
)

// removeSwitchHandler unregisters a robot as a switch
func removeSwitchHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// parse robot ID
		robotID, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			httpWriteError(w, ErrorInvalidRobotID, http.StatusBadRequest)
			return
		}

		// insert into database
		result, err := db.Exec("DELETE FROM switches WHERE robotId = $1", robotID)
		if err != nil {
			log.Printf("Failed to insert user into database: %v", err)
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		count, err := result.RowsAffected()
		if err != nil {
			httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if count == 0 {
			httpWriteError(w, ErrorRobotNotRegistered, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func createRouter(db *sql.DB) *httprouter.Router {
	router := httprouter.New()

	// register routes
	router.POST("/:id", addSwitchHandler(db))
	router.DELETE("/:id", removeSwitchHandler(db))

	return router
}

// logRequests logs the method and URL of each request
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func main() {
	// connect to database
	db := getDb()
	defer db.Close()

	// test connection
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %+v\n", db.Stats())

	// start server
	if err := http.ListenAndServe(":80", logRequests(createRouter(db))); err != nil {
		panic(err)
	}
}
