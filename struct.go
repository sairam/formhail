package main

// NOTE: Domain confirmation will happen through email only by requesting
import (
	"net/mail"
	"net/url"
)

// SingleFormConfig has configuration for a single form
type SingleFormConfig struct {
	Name          string
	UID           string // UID is an alias to the Email, but a randomly generated string
	Email         *mail.Address
	URL           *url.URL // A page means a single page is supported
	URLType       string   // URLType can be page or domain or regexp for URL that needs to be matched
	Confirmed     string   // true / false / spam
	ConfirmedDate string   // datetime at which this confirmation was made

	// Counters to track for incoming
	Counter // TODO incoming counter should be at Domain or Email level instead of form level

	accType *AccountType // Links to an Account Type

	// All notifications to external points can be configured through this
	// Limits apply based on AccountType
	notifications []*Notifier
}

// Notifier is always an outgoing notification sent
type Notifier struct {
	Settings     map[string]string // any other setting like header etc.,
	EndPointURL  string            // https://.... or user@example.com
	EndPointType string            // slack,email,webhook
	Verified     bool              // email requires verification while slack/webhooks don't
	Internal     bool              // internal is an explicit one based on the registered email

	Counter
}

// Counter to track no. of requests processed till ChangeTime. links to AccountLimit through AccountType
type Counter struct {
	Count      int   // current no of requests served
	ChangeTime int64 // Next ChangeTime calculated when Count reaches the Limit
}

// AccountType has a name, description and limits based on the type of channel
type AccountType struct {
	Name   string                  // Basic
	Limits map[string]AccountLimit // Has different Configuration
}

// AccountLimit defines how many requests can be accepted per a period
type AccountLimit struct {
	Type string // incoming, outgoing:slack, outgoing:email, outgoing:webhook
	// Limit & Period are configurable at Account / User level
	// if limit is -1, unlimited will be sent.
	Limit  int // no. of Requests to limit to until ChangeTime
	Period int // no. of seconds from ChangeTime it will reset to ChangeTime += Period & Count = 0
}

// UserSignInRequest is filled up when a user requests a validation/login
type UserSignInRequest struct {
	Email     string
	Domain    *url.URL // generic login / domain related login
	RandomID  string
	Status    string // used / spam / notused
	ReqTime   int64  // requested time request epoch
	ValidTime int64  // valid time of RandomID uses time request epoch
	SEndTime  int64  // Session End Time request epoch
}

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

// ProcessedRequest links incoming request and the config initially provided by the user
type ProcessedRequest struct {
	*IncomingRequest
	*SingleFormConfig
}
