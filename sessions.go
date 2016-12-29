package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/sessions"
	"github.com/sairam/kinli"
)

type session struct{}

func init() {
	kinli.SessionStore = sessions.NewFilesystemStore("./sessions", []byte(os.Getenv("SESSION_STORE")))
	kinli.SessionName = config.SessionName
	kinli.IsAuthed = isAuthed
}

func isAuthed(hc *kinli.HttpContext) bool {
	u := hc.GetSessionData("user")
	if user, ok := u.(*UserSession); ok && user.Email != "" {
		return true
	}
	return false
}

// Login Action called on GET/POST and any error on email adress renders a simple login form
func (session) Login(w http.ResponseWriter, r *http.Request) {
	hc := &kinli.HttpContext{W: w, R: r}

	if r.Method == http.MethodPost {
		err := makeAToken(r)
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
func (session) LogMeIn(w http.ResponseWriter, r *http.Request) {
	loginRequestType := getFirstValue(r.URL.Query(), "type")
	token := getFirstValue(r.URL.Query(), "token")
	report := getFirstValue(r.URL.Query(), "report")
	if loginRequestType != sirequestTypeConfirm && loginRequestType != sirequestTypeLogin {
		http.NotFound(w, r)
		return
	}
	// var page kinli.Page
	if report == "spam" {
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

// UserSession information
type UserSession struct {
	Email  string
	Domain string
}

func getFirstValue(q url.Values, key string) string {
	if len(q[key]) == 0 {
		return ""
	}
	return q[key][0]
}

func (session) Logout(w http.ResponseWriter, r *http.Request) {
	hc := &kinli.HttpContext{W: w, R: r}
	hc.ClearSession()
	hc.RedirectUnlessAuthed("")
}
