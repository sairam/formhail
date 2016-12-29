package controller

import (
	"log"
	"net/http"
	"time"

	"../common"
	"../service"

	"github.com/gorilla/mux"
	"github.com/sairam/kinli"
)

// InitRouter start the router/subrouters
func InitRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/login", service.Session{}.Login)
	r.HandleFunc("/logmein", service.Session{}.LogMeIn).Methods("GET")
	r.HandleFunc("/logout", service.Session{}.Logout).Methods("GET")

	r.HandleFunc("/faq", func(w http.ResponseWriter, r *http.Request) {
		hc := &kinli.HttpContext{W: w, R: r}
		page := kinli.NewPage(hc, "Frequently Asked Questions", "", "", nil)
		kinli.DisplayPage(w, "faq", page)
	}).Methods("GET")

	r.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
		hc := &kinli.HttpContext{W: w, R: r}
		page := kinli.NewPage(hc, "Example Form", "", "", nil)
		kinli.DisplayPage(w, "example", page)
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hc := &kinli.HttpContext{W: w, R: r}
		page := kinli.NewPage(hc, "hello page", "", "", nil)
		kinli.DisplayPage(w, "home", page)
	}).Methods("GET")

	r.HandleFunc("/{uid}", service.FormSubmissionRequest).Methods("POST")

	initStatic(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         common.Config.LocalServer,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}
