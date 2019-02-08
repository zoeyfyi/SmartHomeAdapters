package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

// Need to a) check is user is a valid user
// b) check if robot is registered already
// c) insert into database
// d) takes in a serial, nickname and type, return done/err

var (
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
	url      = os.Getenv("DB_URL")
)
var signingKey = os.Getenv("JWT_SIGNING_KEY")

const (
	ErrorInvalidToken = "Authorization token is invalid"
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

type restError struct {
	Error string `json:"error"`
}

func httpWriteError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(restError{msg})
}

func registerRobotHandler(db *sql.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		robotId := ps.ByName("robotId")
		nickname := ps.ByName("nickname")
		robotType := ps.ByName("robotType")

		log.Printf("/registerRobot/%s/%s/%s", robotId, nickname, robotType)

		// Firstly, check if valid user

		tokenString := r.Header.Get("token")

		// Debugging
		// tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDE5LTAyLTA5VDEzOjI0OjA3LjM0NTE1MjNaIiwiaWQiOiIxIn0.GSaooyU2DPApvJdDesN-4L6TOIuByTglqoMz15b0Qbc"

		if tokenString == "" {
			log.Printf("User did not login")
			httpWriteError(w, "Please login first.", http.StatusBadRequest)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})
		if err != nil {
			log.Printf("Couldnt parse token: %v", err)
			httpWriteError(w, ErrorInvalidToken, http.StatusBadRequest)
			return
		}

		// get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Printf("Could not get token claims")
			httpWriteError(w, ErrorInvalidToken, http.StatusBadRequest)
			return
		}

		// get ID from claims
		id, ok := claims["id"].(string)
		if !ok {
			log.Printf("Token did not have ID claim, claims: %+v", claims)
			httpWriteError(w, ErrorInvalidToken, http.StatusBadRequest)
			return
		}

		// Insert robot into DB with corresponding user ID

		if robotType == "toggle" {
			rows, err := db.Query("SELECT * FROM toggleRobots WHERE serial = $1", robotId)
			if err != nil {
				log.Println("Failed to search database for robot.")
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			for rows.Next() {
				httpWriteError(w, "Robot already exists.", http.StatusInternalServerError)
			}

			_, err = db.Query("INSERT INTO toggleRobots (serial, nickname, robotType, registeredUserId) VALUES ($1, $2, $3, $4)", robotId, nickname, robotType, id)
			if err != nil {
				log.Println("Failed to insert toggle robot into database: %v", err)
			}

		} else if robotType == "range" {
			rows, err := db.Query("SELECT * FROM rangeRobots WHERE serial = $1", robotId)
			if err != nil {
				log.Println("Failed to search database for robot.")
				httpWriteError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			for rows.Next() {
				httpWriteError(w, "Robot already exists.", http.StatusInternalServerError)
			}
			_, err = db.Query("INSERT INTO rangeRobots (serial, nickname, robotType, minimum, maximum, registeredUserId) VALUES ($1, $2, $3, $4, $5, $6)", robotId, nickname, robotType, 0, 0, id)
			if err != nil {
				log.Println("Failed to insert range robot into database: %v", err)
			}
		} else {
			httpWriteError(w, "Error. Invalid robot type.", http.StatusInternalServerError)
		}

	})
}

func createRouter(db *sql.DB) *httprouter.Router {
	router := httprouter.New()
	router.GET("/registerRobot/:robotId/:nickname/:robotType", registerRobotHandler(db))
	return router
}

func main() {
	// register routes
	// mux := http.NewServeMux()

	db := getDb()
	defer db.Close()

	// test database
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Printf("Connected to database: %+v\n", db.Stats())

	// start server
	if err := http.ListenAndServe(":80", createRouter(db)); err != nil {
		panic(err)
	}
}
