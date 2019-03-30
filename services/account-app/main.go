package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"google.golang.org/grpc"

	"github.com/julienschmidt/httprouter"
	"github.com/mrbenshef/SmartHomeAdapters/microservice/userserver"
)

var (
	userserverClient userserver.UserServerClient
)

var (
	loginTemplate    = template.Must(template.ParseFiles("/app/static/login.html"))
	registerTemplate = template.Must(template.ParseFiles("/app/static/register.html"))
	consentTemplate  = template.Must(template.ParseFiles("/app/static/consent.html"))
)

type hydraLoginRequest struct {
	Skip bool `json:"skip"` // whether to skip login or not
	// https://www.ory.sh/docs/hydra/sdk/api#get-an-login-request
}

type hydraLoginAccept struct {
	RedirectTo string `json:"redirect_to"` // url to redirect to
}

type loginTemplateData struct {
	Error error
}

type hydraConsentAccept struct {
	RedirectTo string `json:"redirect_to"` // url to redirect to
}

type consentTemplateData struct {
	Error error
}

type registerTemplateData struct {
	Error   error
	Success bool
}

type LoginRequest struct {
	Remember bool   `json:"remember"`
	Subject  string `json:"subject"`
}

func loginError(w http.ResponseWriter, err error) {
	err = loginTemplate.Execute(w, loginTemplateData{
		Error: err,
	})
	if err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func acceptLogin(w http.ResponseWriter, r *http.Request, email string, password string, challenge string) {
	acceptLoginURL := fmt.Sprintf("https://hydra.halspals.co.uk/oauth2/auth/requests/login/%s/accept", challenge)

	// parse the email and password from the form and attempt to log the user in
	err := r.ParseForm()
	if err != nil {
		log.Printf("failed to parse form: %v", err)
		loginError(w, errors.New("internal error"))
		return
	}

	log.Printf("logging in user with email: %s", email)
	token, err := userserverClient.Login(context.Background(), &userserver.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Printf("bad login: %v", err)
		loginError(w, err)
		return
	}

	// TODO: remove old token handling
	id, err := userserverClient.Authorize(context.Background(), token)
	if err != nil {
		log.Printf("Bad authorize: %v", err)
		loginError(w, errors.New("internal error"))
		return
	}

	// respond to hydra with user ID
	jsonData := &LoginRequest{Remember: false, Subject: id.Id}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(jsonData)
	if err != nil {
		log.Printf("failed to encode hydra responce: %v", err)
		loginError(w, errors.New("internal error"))
		return
	}

	req, err := http.NewRequest("PUT", acceptLoginURL, buf)
	if err != nil {
		log.Printf("error making hydra request: %v", err)
		loginError(w, errors.New("internal error"))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error doing hydra request: %v", err)
		loginError(w, errors.New("internal error"))
		return
	}

	if resp.Body == nil {
		log.Println("received nil responce body from hydra")
		loginError(w, errors.New("internal error"))
		return
	}

	var loginAccept hydraLoginAccept
	err = json.NewDecoder(resp.Body).Decode(&loginAccept)
	if err != nil {
		log.Printf("error reading hydra login accept body: %v", err)
		loginError(w, errors.New("internal error"))
		return
	}

	// redirect back to hydra
	http.Redirect(w, r, loginAccept.RedirectTo, http.StatusMovedPermanently)
}

func postLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// parse post form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	challenge := r.Form.Get("challenge")

	acceptLogin(w, r, email, password, challenge)
}

func getLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("template: %v", loginTemplate)

	// get hydra challenge
	challenge := r.URL.Query().Get("login_challenge")
	if challenge == "" {
		log.Println("no login challenge")
		loginError(w, errors.New("bad OAuth client"))
		return
	}

	// get infomation about the login request from hydra
	urlGet := fmt.Sprintf("https://hydra.halspals.co.uk/oauth2/auth/requests/login/%s", challenge)
	resp, err := http.DefaultClient.Get(urlGet)
	if err != nil {
		// problem with hydra, fallback to login form
		log.Printf("error doing hydra request: %v", err)
		loginError(w, errors.New("internal error"))
		return
	}

	// check responce is valid
	if resp.StatusCode != 200 || resp.Body == nil {
		log.Printf("invalid responce from hydra (%d)", resp.StatusCode)
		loginError(w, errors.New("internal error"))
		return
	}

	// decode hydra responce
	var loginRequest hydraLoginRequest
	err = json.NewDecoder(resp.Body).Decode(&loginRequest)
	if err != nil {
		log.Printf("failed to decode hydra responce, error: %v", err)

		err = loginTemplate.Execute(w, loginTemplateData{
			Error: errors.New("internal error"),
		})
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
		return
	}

	// hydra recognizes user, skip login
	if loginRequest.Skip {
		log.Println("skipping login")
		loginError(w, errors.New("login skip is unimplemented"))
		// TODO: need to get info from userserver without password
		// acceptLogin(w, r, challenge, "", "")
		return
	}

	// render login form
	err = loginTemplate.Execute(w, loginTemplateData{})
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func registerError(w http.ResponseWriter, err error) {
	err = registerTemplate.Execute(w, registerTemplateData{
		Error: err,
	})
	if err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func postRegisterHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// parse post form
	err := r.ParseForm()
	if err != nil {
		registerError(w, errors.New("internal error"))
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	confirmPassword := r.PostForm.Get("password_confirm")

	// check passwords match
	if password != confirmPassword {
		registerError(w, errors.New("passwords do not match"))
		return
	}

	// register user
	_, err = userserverClient.Register(context.Background(), &userserver.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		registerError(w, err)
		return
	}

	// return success
	err = registerTemplate.Execute(w, registerTemplateData{
		Success: true,
	})
	if err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func getRegisterHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := registerTemplate.Execute(w, registerTemplateData{})
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func consentError(w http.ResponseWriter, err error) {
	err = consentTemplate.Execute(w, consentTemplateData{
		Error: err,
	})
	if err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func acceptConsent(w http.ResponseWriter, r *http.Request, challenge string) {
	acceptConsentURL := fmt.Sprintf("https://hydra.halspals.co.uk/oauth2/auth/requests/consent/%s/accept", challenge)

	jsonData := []byte(`{ "remember":false, "grant_scope":["openid"]}`)
	req, err := http.NewRequest("PUT", acceptConsentURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("error making hydra request: %v", err)
		consentError(w, errors.New("internal error"))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error doing hydra request: %v", err)
		consentError(w, errors.New("internal error"))
		return
	}

	if resp.Body == nil {
		log.Println("recevied nil body from hydra")
		consentError(w, errors.New("internal error"))
		return
	}

	var consentAccept hydraConsentAccept
	err = json.NewDecoder(resp.Body).Decode(&consentAccept)
	if err != nil {
		log.Printf("failed to decode JSON: %v", err)
		consentError(w, errors.New("internal error"))
		return
	}

	// redirect back to hydra
	http.Redirect(w, r, consentAccept.RedirectTo, http.StatusMovedPermanently)
}

func postConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get hydra challenge
	err := r.ParseForm()
	if err != nil {
		log.Printf("failed to parse form: %v", err)
		consentError(w, errors.New("bad request"))
		return
	}

	challenge := r.Form.Get("challenge")
	acceptConsent(w, r, challenge)
}

func getConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := consentTemplate.Execute(w, consentTemplateData{
		Error: nil,
	})
	if err != nil {
		log.Printf("error rendering template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	// TODO: resolve this
	// nolint
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	router := httprouter.New()
	router.GET("/login", getLoginHandler)
	router.POST("/login", postLoginHandler)
	router.GET("/register", getRegisterHandler)
	router.POST("/register", postRegisterHandler)
	router.GET("/consent", getConsentHandler)
	router.POST("/consent", postConsentHandler)
	router.Handler("GET", "/static/*filepath", http.StripPrefix("/static/", http.FileServer(http.Dir("/app/static/"))))

	userserverConn, err := grpc.Dial("userserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer userserverConn.Close()
	userserverClient = userserver.NewUserServerClient(userserverConn)
	// start server
	if err := http.ListenAndServe(":80", router); err != nil {
		panic(err)
	}
}
