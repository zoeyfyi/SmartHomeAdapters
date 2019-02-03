package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

var client = http.DefaultClient

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
	RobotID      int  `json:"robotId"`
	IsOn         bool `json:"isOn"`
	IsCalibrated bool `json:"isCalibrated"`
	OnAngle      int  `json:"onAngle"`
	OffAngle     int  `json:"offAngle"`
	RestAngle    int  `json:"restAngle"`
}

type restError struct {
	Error string `json:"error"`
}

func httpWriteError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(restError{msg})
}

// error messages for add switch route
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

		switchRobot := switchRobot{
			RobotID: robotID,
			IsOn:    *addSwitchBody.IsOn,
		}

		// insert into database
		_, err = db.Exec(
			"INSERT INTO switches(robotId, isOn, onAngle, offAngle, restAngle, isCalibrated) VALUES($1, $2, $3, $4, $5, $6)",
			switchRobot.RobotID,
			switchRobot.IsOn,
			switchRobot.OnAngle,
			switchRobot.OffAngle,
			switchRobot.RestAngle,
			switchRobot.IsCalibrated,
		)
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
		json.NewEncoder(w).Encode(switchRobot)
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

func setServo(w http.ResponseWriter, angle int) bool {
	// send servo commands
	url := fmt.Sprintf("robotserver/servo/%d", angle)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		w.Write(buf.Bytes())
		return false
	}

	return true
}

func setIsOn(w http.ResponseWriter, db *sql.DB, robotID int, isOn bool) bool {
	res, err := db.Exec("UPDATE switches SET isOn = $1 WHERE robotId = $2", isOn, robotID)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Failed to get the amount of rows affected: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	if rowsAffected != 1 {
		log.Printf("Expected to update exactly 1 row, rows updated: %d\n", rowsAffected)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	return true
}

// errors for light on/off routes
const (
	ErrorNotCalibrated = "Robot has not been calibrated"
	ErrorSwitchOn      = "Switch is already on (use force to override)"
	ErrorSwitchOff     = "Switch is already off (use force to override)"
)

func setSwitch(w http.ResponseWriter, db *sql.DB, robotID int, setOn bool, force bool) {
	// get switch status from database
	var (
		rID          int
		isOn         bool
		isCalibrated bool
		onAngle      int
		offAngle     int
		restAngle    int
	)

	row := db.QueryRow("SELECT robotId, isOn, isCalibrated, onAngle, offAngle, restAngle FROM switches WHERE robotId = $1", robotID)
	err := row.Scan(&rID, &isOn, &isCalibrated, &onAngle, &offAngle, &restAngle)
	if err != nil {
		log.Printf("Failed to scan database: %v", err)
		httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// check if it has been calibrated
	if !isCalibrated {
		httpWriteError(w, ErrorNotCalibrated, http.StatusBadRequest)
		return
	}

	// check if the light is already on
	if !force {
		// check if switch is already on/off
		if isOn == setOn {
			if isOn {
				httpWriteError(w, ErrorSwitchOn, http.StatusBadRequest)
			} else {
				httpWriteError(w, ErrorSwitchOff, http.StatusBadRequest)
			}
			return
		}
	}

	targetAngle := 0
	if setOn {
		targetAngle = onAngle
	} else {
		targetAngle = offAngle
	}

	// send servo commands
	if ok := setServo(w, targetAngle); !ok {
		return
	}
	time.Sleep(time.Second)
	if ok := setServo(w, restAngle); !ok {
		return
	}

	// update database
	if ok := setIsOn(w, db, robotID, true); !ok {
		return
	}

	w.WriteHeader(http.StatusOK)
}

type onSwitchBody struct {
	// Force determines weather a servo command is sent even if the robot is already on
	Force bool `json:"force"`
}

func onHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// parse robot ID
		robotID, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			httpWriteError(w, ErrorInvalidRobotID, http.StatusBadRequest)
			return
		}

		var onSwitchBody onSwitchBody

		if r.Body != nil {
			// decode body
			err = json.NewDecoder(r.Body).Decode(&onSwitchBody)
			if err != nil {
				httpWriteError(w, ErrorInvalidJSON, http.StatusBadRequest)
				return
			}
		}

		// update switch
		setSwitch(w, db, robotID, true, onSwitchBody.Force)
	}
}

type offSwitchBody struct {
	// Force determines weather a servo command is sent even if the robot is already on
	Force bool `json:"force"`
}

func offHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// parse robot ID
		robotID, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			httpWriteError(w, ErrorInvalidRobotID, http.StatusBadRequest)
			return
		}

		// decode body
		var offSwitchBody offSwitchBody
		if r.Body != nil {
			err = json.NewDecoder(r.Body).Decode(&offSwitchBody)
			if err != nil {
				httpWriteError(w, ErrorInvalidJSON, http.StatusBadRequest)
				return
			}
		}

		// update switch
		setSwitch(w, db, robotID, false, offSwitchBody.Force)

		// update database
		setIsOn(w, db, robotID, false)
	}
}

func createRouter(db *sql.DB) *httprouter.Router {
	router := httprouter.New()

	// register routes
	router.POST("/:id", addSwitchHandler(db))
	router.DELETE("/:id", removeSwitchHandler(db))
	router.PATCH("/:id/on", onHandler(db))
	router.PATCH("/:id/off", offHandler(db))

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
