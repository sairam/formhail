package service

import (
	"log"
	"net/http"

	"../helper"
	"github.com/sairam/kinli"
)

// Session handles are login/logout related business logic
type Session struct{}

const (
	RequestActionLogin   = "login"
	RequestActionConfirm = "confirm"
	ReportActionSpam     = "spam"
)

// IsUserAuthed is a helper method to identify if user is logged in
func IsUserAuthed(hc *kinli.HttpContext) bool {
	u := hc.GetSessionData("user")
	if user, ok := u.(*UserSession); ok && user.Email != "" {
		return true
	}
	return false
}

// Login Action called on GET/POST and any error on email adress renders a simple login form
func (Session) Login(w http.ResponseWriter, r *http.Request) {
	hc := &kinli.HttpContext{W: w, R: r}

	if r.Method == http.MethodPost {
		// TODO - requestToAuthenticate should create a domain/email combo if not present and sent the confirm link
		err := requestToAuthenticate(r)
		if err != nil {
			hc.AddFlash("There was a problem processing your request. A human slave shall look into this and will get back to you")
			// TODO send the request we got with debug information to an error stream
			log.Println(err)
		} else {
			hc.AddFlash("You should receive an email momentarily")
		}
		// display only success message
	} else {
		// render html for login request with optional captcha
		// Process the request by check form parameters
	}
	page := kinli.NewPage(hc, "Login", "", "", nil)
	kinli.DisplayPage(w, "login", page)
}

// LogMeIn logs in the user from the link sent to the user in email
// redirects to loggedin Page is successful or to /faq#errors page
func (Session) LogMeIn(w http.ResponseWriter, r *http.Request) {
	loginRequestType := helper.GetFirstValue(r.URL.Query(), "type")
	token := helper.GetFirstValue(r.URL.Query(), "token")
	report := helper.GetFirstValue(r.URL.Query(), "report")

	if loginRequestType != RequestActionLogin && loginRequestType != RequestActionConfirm {
		http.NotFound(w, r)
		return
	}
	// var page kinli.Page
	if report == ReportActionSpam {
		log.Println("spammmed")
		// report spam
		// set page
		// page = kinli.NewPage(hc, "Login", "", "", nil)
		// kinli.DisplayPage(w, "login", page)
		return
	}

	var err error
	var user *UserSession
	if loginRequestType == RequestActionLogin {
		user, err = verifyLoginToken(token)
	} else if loginRequestType == RequestActionConfirm {
		user, err = verifyUserConfirmToken(token)
	}

	hc := &kinli.HttpContext{W: w, R: r}
	if err != nil {
		hc.AddFlash(err.Error())
		http.NotFound(w, r)
		return
	}

	hc.SetSessionData("user", user)
	hc.RedirectAfterAuth("Logged in")
}

// Logout logs the user out by clearing the session
func (Session) Logout(w http.ResponseWriter, r *http.Request) {
	hc := &kinli.HttpContext{W: w, R: r}
	hc.ClearSession()
	hc.RedirectUnlessAuthed("Kindly Login")
}
