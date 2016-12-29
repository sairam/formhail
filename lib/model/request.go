package model

import (
	"net/mail"
	"net/url"
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
