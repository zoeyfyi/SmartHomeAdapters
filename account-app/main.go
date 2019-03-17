package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/mrbenshef/SmartHomeAdapters/userserver/userserver"
)

var (
	userserverClient userserver.UserServerClient
)
var (
	loginTemplate   = template.Must(template.ParseFiles("/app/static/login.html"))
	consentTemplate = template.Must(template.ParseFiles("/app/static/consent.html"))
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

type hydraConsentRequest struct {
	Skip bool `json:"skip"`
}

type hydraConsentAccept struct {
	RedirectTo string `json:"redirect_to"` // url to redirect to
}

type consentTemplateData struct {
	Error error
}

type LoginRequest struct {
	Remember bool   `json:"remember"`
	Subject  string `json:"subject"`
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func acceptLogin(w http.ResponseWriter, r *http.Request, challenge string) {
	urlPut := fmt.Sprintf("https://hydra.halspals.co.uk/oauth2/auth/requests/login/%s/accept", challenge)

	// parse the email and password from the form and attempt to log the user in
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	log.Printf("Logging in user with email: %s", email)
	token, err := userserverClient.Login(context.Background(), &userserver.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Bad login: %v", err)
	}

	id, err := userserverClient.Authorize(context.Background(), token)

	if err != nil {
		log.Fatalf("Bad authorize: %v", err)
	}

	jsonData := &LoginRequest{Remember: false, Subject: id.Id}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(jsonData)

	data := url.Values{}
	data.Set("\"remember\"", "false")
	data.Set("\"subject\"", "\"subject123\"")

	req, err := http.NewRequest("PUT", urlPut, buf)
	if err != nil {
		loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("Internal error"),
		})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("Internal error"),
		})
		return
	}

	if resp.Body == nil {
		loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("Internal error"),
		})
		return
	}
	var loginAccept hydraLoginAccept

	body, err := ioutil.ReadAll(resp.Body)

	reader := bytes.NewReader(body)
	err = json.NewDecoder(reader).Decode(&loginAccept)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	strings.Replace(loginAccept.RedirectTo, "\\u0026", "&", -1)

	// redirect back to hydra
	http.Redirect(w, r, loginAccept.RedirectTo, http.StatusMovedPermanently)

}

func postLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// parse post form
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	challenge := r.Form.Get("challenge")

	// check email/password
	if email == "" {
		loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("Email is blank"),
		})
		return
	}
	if password == "" {
		loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("Password is blank"),
		})
		return
	}

	acceptLogin(w, r, challenge)
}

func getLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("template: %v", loginTemplate)

	// get hydra challenge
	challenge := r.URL.Query().Get("login_challenge")

	// get infomation about the login request from hydra
	urlGet := fmt.Sprintf("https://hydra.halspals.co.uk/oauth2/auth/requests/login/%s", challenge)
	resp, err := http.DefaultClient.Get(urlGet)
	if err != nil {
		loginTemplate.Execute(w, loginTemplateData{
			Error: err,
		})
		return
	}

	if resp.Body != nil {
		var loginRequest hydraLoginRequest
		json.NewDecoder(resp.Body).Decode(&loginRequest)

		if loginRequest.Skip {
			acceptLogin(w, r, challenge)
			return
		}
	}

	err = loginTemplate.Execute(w, loginTemplateData{
		Error: nil,
	})
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func acceptConsent(w http.ResponseWriter, r *http.Request, challenge string) {
	urlPut := fmt.Sprintf("https://hydra.halspals.co.uk/oauth2/auth/requests/consent/%s/accept", challenge)

	jsonData := []byte(`{ "remember":false, "grant_scope":["openid"]}`)

	req, err := http.NewRequest("PUT", urlPut, bytes.NewBuffer(jsonData))
	if err != nil {
		consentTemplate.Execute(w, consentTemplateData{
			Error: fmt.Errorf("Internal error"),
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		consentTemplate.Execute(w, consentTemplateData{
			Error: fmt.Errorf("Internal error"),
		})
		return
	}

	if resp.Body == nil {
		consentTemplate.Execute(w, consentTemplateData{
			Error: fmt.Errorf("Internal error"),
		})
		return
	}

	var consentAccept hydraConsentAccept
	json.NewDecoder(resp.Body).Decode(&consentAccept)

	// redirect back to hydra
	http.Redirect(w, r, consentAccept.RedirectTo, http.StatusMovedPermanently)
}

func postConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get hydra challenge
	r.ParseForm()
	challenge := r.Form.Get("challenge")

	acceptConsent(w, r, challenge)
}

func getConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := consentTemplate.Execute(w, consentTemplateData{
		Error: nil,
	})
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	router := httprouter.New()
	router.GET("/login", getLoginHandler)
	router.POST("/login", postLoginHandler)
	router.GET("/consent", getConsentHandler)
	router.POST("/consent", postConsentHandler)

	userserverConn, err := grpc.Dial("userserver:80", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer userserverConn.Close()
	userserverClient = userserver.NewUserServerClient(userserverConn)
	// start server
	if err := http.ListenAndServe(":4001", router); err != nil {
		panic(err)
	}
}
