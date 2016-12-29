package model

import (
	"errors"
	"log"
	"net/mail"
	"net/url"
	"strings"

	"../helper"
)

const (
	ccListLimit = 10
)

// IncomingRequest is the incoming structure to fill when a form is submitted
type IncomingRequest struct {
	Referral   *url.URL        // mandatory to be verified
	Identifier string          // Identifier is the email or UID present in the form POST url
	IDType     string          // type is email or id
	ReplyTo    *mail.Address   // optional
	NextPage   *url.URL        // optional
	Subject    string          // optional
	Cc         []*mail.Address // optional
	Format     []string        // optional, default html , set to plain
	// Gotcha     string          // should be ignored when set to any string other than blank

	Message map[string][]string // url.Values from the form after removing the optional ones

	DateTime   int64 // datetime at which we have received the request
	RemoteAddr string
}

type FormRequest struct {
	Referral   *url.URL `schema:"-"`
	Identifier string   `schema:"-"`
	ReplyTo    string   `schema:"_replyto"`
	NextPage   string   `schema:"_next"`
	Subject    string   `schema:"_subject"`
	CcString   string   `schema:"_cc"`
	Format     string   `schema:"_format"`
	Gotcha     string   `schema:"_gotcha"`
}

// Validate by querying the data store
func (not *IncomingRequest) Validate() error {
	return nil
}

// ParseFormFields parses the input data and fills in the struct
func (not *IncomingRequest) ParseFormFields(referrer, requestURI string, form *FormRequest) error {

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

	requestID := strings.Trim(uri.Path, "/")
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
		ccList, err = mail.ParseAddressList(strings.Trim(form.CcString, ","))
		if err == nil {
			maxCount := helper.Min(len(ccList), ccListLimit)
			ccList = ccList[:maxCount]
			not.Cc = ccList
		} else {
			log.Println(err)
		}
	} else {
		not.Cc = make([]*mail.Address, 0)
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

	if form.Format == "plain" {
		not.Format = []string{"plain"}
	} else {
		not.Format = []string{"html", "plain"}
	}

	if addr, err := mail.ParseAddress(form.ReplyTo); err == nil {
		not.ReplyTo = addr
	}

	return nil
}
