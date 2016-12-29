package service

import (
	"log"
	"net/http"

	"../helper"
	"github.com/sairam/kinli"
)

// IsUserAuthed is a helper method to identify if user is logged in
func IsUserAuthed(hc *kinli.HttpContext) bool {
	u := hc.GetSessionData("user")
	if user, ok := u.(*UserSession); ok && user.Email != "" {
		return true
	}
	return false
}

// Session handles are login/logout related business logic
type Session struct{}

// Login Action called on GET/POST and any error on email adress renders a simple login form
func (Session) Login(w http.ResponseWriter, r *http.Request) {
	hc := &kinli.HttpContext{W: w, R: r}

	if r.Method == http.MethodPost {
		err := requestToAuthenticate(r)
		if err != nil {
			log.Println(err)
			// page . display error message
		}
		// display only success message
	} else {
		// render html for login request with optional captcha
		// Process the request by check form parameters
	}
	page := kinli.NewPage(hc, "Login", "", "", nil)
	kinli.DisplayPage(w, "login", page)
}

// GET url user uses from the email
// redirects to loggedin Page is successful or to /faq#errors page
func (Session) LogMeIn(w http.ResponseWriter, r *http.Request) {
	loginRequestType := helper.GetFirstValue(r.URL.Query(), "type")
	token := helper.GetFirstValue(r.URL.Query(), "token")
	report := helper.GetFirstValue(r.URL.Query(), "report")
	if loginRequestType != "confirm" && loginRequestType != "login" {
		http.NotFound(w, r)
		return
	}
	// var page kinli.Page
	if report == "spam" {
		log.Println("spammmed")
		// report spam
		// set page
		// page = kinli.NewPage(hc, "Login", "", "", nil)
		// kinli.DisplayPage(w, "login", page)
		return
	}

	var err error
	var user *UserSession
	if loginRequestType == "login" {
		user, err = verifyLoginToken(token)
	} else if loginRequestType == "confirm" {
		user, err = verifyUserConfirmToken(token)
	}
	if err != nil {
		log.Println(err)
		return
		// display page with text as error
	}
	hc := &kinli.HttpContext{W: w, R: r}
	hc.SetSessionData("user", user)
	hc.RedirectAfterAuth("Logged in")

	// user is redirected to / or /dashboard based on the result
}

func (Session) Logout(w http.ResponseWriter, r *http.Request) {
	hc := &kinli.HttpContext{W: w, R: r}
	hc.ClearSession()
	hc.RedirectUnlessAuthed("")
}
