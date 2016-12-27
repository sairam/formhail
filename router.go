package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sairam/kinli"
)

func initRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/login", session{}.Login)
	r.HandleFunc("/logmein", session{}.LogMeIn).Methods("GET")
	r.HandleFunc("/logout", session{}.Logout).Methods("GET")

	r.HandleFunc("/faq", func(w http.ResponseWriter, r *http.Request) {
		hc := &kinli.HttpContext{W: w, R: r}
		page := kinli.NewPage(hc, "Frequently Asked Questions", "", "", nil)
		kinli.DisplayPage(w, "faq", page)
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hc := &kinli.HttpContext{W: w, R: r}
		page := kinli.NewPage(hc, "hello page", "", "", nil)
		kinli.DisplayPage(w, "home", page)
	}).Methods("GET")

	r.HandleFunc("/", NewSubmissionRequest).Methods("POST")

	initStatic(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:3000",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}
