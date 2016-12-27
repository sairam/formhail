package main

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/sairam/kinli"
)

type session struct{}

func init() {
	kinli.SessionStore = sessions.NewFilesystemStore("./sessions", []byte(os.Getenv("SESSION_STORE")))
	kinli.SessionName = "_formr"
	kinli.IsAuthed = isAuthed
}

func isAuthed(hc *kinli.HttpContext) bool {
	return true
}

// Login Action called on POST and any error on email adress renders a simple login form
func (session) Login(w http.ResponseWriter, r *http.Request) {
	hc := &kinli.HttpContext{W: w, R: r}

	// if confirmToken="abc". use Login()
	// render html for login request with optional captcha
	// Process the request by check form parameters

	page := kinli.NewPage(hc, "Login", "", "", nil)
	kinli.DisplayPage(w, "login", page)

}

// GET url user uses from the email
// redirects to loggedin Page is successful or to /faq#errors page
func (session) LogMeIn(w http.ResponseWriter, r *http.Request) {
	// process login request on GET
}

func (session) Logout(w http.ResponseWriter, r *http.Request) {
	hc := &kinli.HttpContext{W: w, R: r}
	hc.ClearSession()
	hc.RedirectUnlessAuthed("")
}
