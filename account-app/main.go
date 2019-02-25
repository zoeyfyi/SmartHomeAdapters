package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

var (
	loginTemplate   = template.Must(template.ParseFiles("static/login.html"))
	consentTemplate = template.Must(template.ParseFiles("static/consent.html"))
)

type hydraLoginRequest struct {
	Skip bool `json:"skip"` // weather to skip login or not
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

func acceptLogin(w http.ResponseWriter, r *http.Request, challenge string) {
	url := fmt.Sprintf("http://hydra/oauth/auth/requests/login/%s/accept", challenge)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("Internal error"),
		})
		return
	}

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
	json.NewDecoder(resp.Body).Decode(loginAccept)

	// redirect back to hydra
	http.Redirect(w, r, loginAccept.RedirectTo, http.StatusMovedPermanently)
}

func postLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get hydra challenge
	challenge := r.URL.Query().Get("login_challenge")

	// parse post form
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

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
	url := fmt.Sprintf("http://hydra/oauth2/auth/requests/login/%s", challenge)
	resp, err := http.DefaultClient.Get(url)
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
	url := fmt.Sprintf("http://hydra/oauth/auth/requests/consent/%s/accept", challenge)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		consentTemplate.Execute(w, consentTemplateData{
			Error: fmt.Errorf("Internal error"),
		})
		return
	}

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
	json.NewDecoder(resp.Body).Decode(consentAccept)

	// redirect back to hydra
	http.Redirect(w, r, consentAccept.RedirectTo, http.StatusMovedPermanently)
}

func postConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get hydra challenge
	challenge := r.URL.Query().Get("login_challenge")

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
	router := httprouter.New()
	router.GET("/login", getLoginHandler)
	router.POST("/login", postLoginHandler)
	router.GET("/consent", getConsentHandler)
	router.POST("/consent", postConsentHandler)

	// start server
	if err := http.ListenAndServe(":4001", router); err != nil {
		panic(err)
	}
}
