package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

func acceptLogin(w http.ResponseWriter, r *http.Request, challenge string) {
	urlPut := fmt.Sprintf("https://hydra.halspals.co.uk/oauth2/auth/requests/login/%s/accept", challenge)

	// parse the email and password from the form and attempt to log the user in
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

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
	err = json.NewEncoder(buf).Encode(jsonData)
	if err != nil {
		log.Printf("failed to encode hydra responce: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := url.Values{}
	data.Set("\"remember\"", "false")
	data.Set("\"subject\"", "\"subject123\"")

	req, err := http.NewRequest("PUT", urlPut, buf)
	if err != nil {
		err = loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("internal error"),
		})
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("internal error"),
		})
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
		return
	}

	if resp.Body == nil {
		err = loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("internal error"),
		})
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
		return
	}
	var loginAccept hydraLoginAccept

	err = json.NewDecoder(resp.Body).Decode(&loginAccept)
	if err != nil {
		log.Printf("error reading hydra login accept body: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// TODO: strings.Replace is a pure function, do we need this?
	// strings.Replace(loginAccept.RedirectTo, "\\u0026", "&", -1)

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

	// check email/password
	if email == "" {
		err = loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("email is blank"),
		})
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
		return
	}
	if password == "" {
		err = loginTemplate.Execute(w, loginTemplateData{
			Error: fmt.Errorf("password is blank"),
		})
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
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
		err = loginTemplate.Execute(w, loginTemplateData{
			Error: err,
		})
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
		return
	}

	if resp.Body != nil {
		var loginRequest hydraLoginRequest
		err = json.NewDecoder(resp.Body).Decode(&loginRequest)
		if err != nil {
			http.Error(w, "failed to decode JSON", http.StatusBadRequest)
		}

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

func postRegisterHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// parse post form
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	confirmPassword := r.PostForm.Get("password_confirm")

	if password != confirmPassword {
		err := registerTemplate.Execute(w, registerTemplateData{
			Error: fmt.Errorf("passwords do not match"),
		})
		if err != nil {
			log.Printf("Error rendering template: %v", err)
		}
		return
	}

	_, err := userserverClient.Register(context.Background(), &userserver.RegisterRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		err = registerTemplate.Execute(w, registerTemplateData{
			Error: err,
		})
		if err != nil {
			log.Printf("Error rendering template: %v", err)
		}
	}

	err = registerTemplate.Execute(w, registerTemplateData{
		Success: true,
	})
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func getRegisterHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := registerTemplate.Execute(w, registerTemplateData{})
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func acceptConsent(w http.ResponseWriter, r *http.Request, challenge string) {
	urlPut := fmt.Sprintf("https://hydra.halspals.co.uk/oauth2/auth/requests/consent/%s/accept", challenge)

	jsonData := []byte(`{ "remember":false, "grant_scope":["openid"]}`)

	req, err := http.NewRequest("PUT", urlPut, bytes.NewBuffer(jsonData))
	if err != nil {
		err = consentTemplate.Execute(w, consentTemplateData{
			Error: fmt.Errorf("internal error"),
		})
		if err != nil {
			log.Printf("error rendering template: %v", err)
		}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = consentTemplate.Execute(w, consentTemplateData{
			Error: fmt.Errorf("internal error"),
		})
		if err != nil {
			log.Printf("Error rendering template: %v", err)
		}
		return
	}

	if resp.Body == nil {
		err = consentTemplate.Execute(w, consentTemplateData{
			Error: fmt.Errorf("internal error"),
		})
		if err != nil {
			log.Printf("Error rendering template: %v", err)
		}
		return
	}

	var consentAccept hydraConsentAccept
	err = json.NewDecoder(resp.Body).Decode(&consentAccept)
	if err != nil {
		http.Error(w, "failed to decode JSON", http.StatusBadRequest)
	}

	// redirect back to hydra
	http.Redirect(w, r, consentAccept.RedirectTo, http.StatusMovedPermanently)
}

func postConsentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get hydra challenge
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
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
