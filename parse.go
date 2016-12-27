package main

import (
	"errors"
	"net/http"
	"net/mail"
	"net/url"
	"time"

	"github.com/gorilla/schema"
)

const (
	ccListLimit = 10
)

type formRequest struct {
	Referral   *url.URL `schema:"-"`
	Identifier string   `schema:"-"`
	ReplyTo    string   `schema:"_replyto"`
	NextPage   string   `schema:"_next"`
	Subject    string   `schema:"_subject"`
	CcString   string   `schema:"_cc"`
	Format     string   `schema:"_format"`
	Gotcha     string   `schema:"_gotcha"`
}

// NewSubmissionRequest handles the incoming POST requests
func NewSubmissionRequest(w http.ResponseWriter, r *http.Request) {
	err := formHandler(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) error {
	var err error

	not := &IncomingRequest{}
	not.DateTime = time.Now().Unix()
	not.RemoteAddr = r.RemoteAddr

	err = r.ParseForm()
	if err != nil {
		return errors.New("Error Parsing Form")
	}
	form := new(formRequest)
	decoder := schema.NewDecoder()
	decoder.Decode(form, r.PostForm)

	err = not.ParseFormFields(r.Referer(), r.RequestURI, form)
	if err != nil {
		return err
	}
	not.Message = removeRestrictedFields(r.PostForm)
	return nil
}

// Validate by querying the data store
func (not *IncomingRequest) Validate() error {
	return nil
}

// ParseFormFields parses the input data and fills in the struct
func (not *IncomingRequest) ParseFormFields(referrer, requestURI string, form *formRequest) error {

	var err error

	form.Referral, err = url.Parse(referrer)
	if err != nil {
		return errors.New("Referral is blank")
	}
	not.Referral = form.Referral

	uri, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return errors.New("Request URI is not parsable")
	}

	requestID := uri.Path
	emailID, err := mail.ParseAddress(requestID)
	if err == nil {
		not.Identifier = emailID.Address
		not.IDType = "email"
	} else {
		not.Identifier = requestID
		not.IDType = "requestID"
	}

	if form.Gotcha != "" {
		return errors.New("Gotcha is set. Expected not to exist")
	}

	if form.CcString != "" {
		var ccList []*mail.Address
		ccList, err = mail.ParseAddressList(form.CcString)
		if err == nil {
			maxCount := min(len(ccList), ccListLimit)
			ccList = ccList[:maxCount]
			not.Cc = ccList
		}
	}

	if form.Subject == "" {
		not.Subject = "Form Filled at " + form.Referral.Path
	} else {
		not.Subject = form.Subject
	}

	if form.NextPage == "" {
		not.NextPage = form.Referral
	} else {
		not.NextPage, err = url.Parse(form.NextPage)
		if err != nil {
			not.NextPage = form.Referral
		}
	}

	if form.Format == "text" {
		not.Format = []string{"text"}
	} else {
		not.Format = []string{"html", "text"}
	}

	if addr, err := mail.ParseAddress(form.ReplyTo); err == nil {
		not.ReplyTo = addr
	}

	return nil
}

func removeRestrictedFields(r url.Values) url.Values {
	for _, key := range []string{"_cc", "_replyto", "_next", "subject", "_format"} {
		r.Del(key)
	}
	return r
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
