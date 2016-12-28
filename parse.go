package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
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
	fmt.Println(not)

	// cases - user/domain
	// 1. is not registered
	// 2. is blacklisted / limit exceeded
	// 3. has enough limit

	pr, err := newProcessedRequest(not)
	var message string
	if pr.SingleFormConfig == nil || err != nil {
		log.Println(err)
		sfc, err := newSingleFormConfig(not)
		log.Printf("%+v", sfc)
		if err != nil {
			message = "error. smth"
		}
		message = "First time user, check email"
		// is not registered
	} else if pr.YetToBeConfirmed() {
		fmt.Println("to be confirmed")
	} else if pr.IsBlacklisted() {
		fmt.Println("blacklisted")
	} else if pr.DidLimitReach() {
		fmt.Println("Limit Reached")
	} else {
		pr.IncrIncoming()
		if !pr.save() {
			fmt.Println("save notifications failed")
		} else {
			pr.SendNotifications()
		}
	}
	fmt.Println(message)
	return nil
}

// responsible for creating and sending an email to the user
func newSingleFormConfig(not *IncomingRequest) (*SingleFormConfig, error) {
	if not.IDType != "email" {
		return nil, errors.New("Incoming Request should be an email")
	}
	page := not.Referral.Path
	page = strings.Trim(page, "/")
	if page == "" {
		page = "homepage"
	}
	email, _ := mail.ParseAddress(not.Identifier)
	uri, _ := url.Parse(not.Referral.Host)
	// TODO add err
	sfc := &SingleFormConfig{
		Name:        page,
		UID:         RandString(20),
		Email:       email,
		URL:         uri,
		Confirmed:   formConfigRequested,
		AccountType: accountTypeBasic,
		Counter: Counter{
			Count:      0,
			ChangeTime: 0,
		},
	}
	sfc.ID = sfc.autoincr()
	sfc.save()
	sfc.Index()
	return sfc, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandString helper method
// source: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// RandStringBytesMaskImprSrc
func RandString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// SendNotifications sends notifications to all subscriptions
func (pr *ProcessedRequest) SendNotifications() {
	// take each of the outgoing notifications and run them through
	for _, not := range pr.Notifications {
		if not.EndPointType == endpointTypeEmail {
			go func() { pr.sendEmail() }()
		} else if not.EndPointType == endpointTypeSlack {
			pr.sendToSlack(not.EndPointURL)
		} else if not.EndPointType == endpointTypeWebhook {
			pr.sendToWebhook(not.EndPointURL)
		} else {
			// not supported
		}
	}
}

// YetToBeConfirmed ..
func (c *SingleFormConfig) YetToBeConfirmed() bool {
	return c.Confirmed == formConfigUnconfirmed || c.Confirmed == formConfigRequested
}

// IsBlacklisted ..
func (c *SingleFormConfig) IsBlacklisted() bool {
	return c.Confirmed == formConfigSpam
}

// DidLimitReach checks if we reached the limit for the account?
// checks for incoming requests
// TODO verify the limit based on the account type
// incr with a lock and save to store
func (c *SingleFormConfig) DidLimitReach() bool {
	// no need to change the time
	if c.ChangeTime == 0 {
		c.ChangeTime = time.Now().Unix() - 10
	}
	accLimit := c.accType.Limits["incoming:form"]
	currentTime := time.Now().Unix()
	for currentTime > c.ChangeTime {
		fmt.Println("time diff is ", currentTime-c.ChangeTime)
		c.ChangeTime += accLimit.Period
		c.Count = 0
	}
	if c.Count < accLimit.Limit {
		return false
	}
	return true
}

// IncrIncoming usage field
// TODO incr with a lock
func (c *SingleFormConfig) IncrIncoming() {
	c.Count = c.Count + 1
}

// should be done with lock to avoid concurrent transactions
// looks through all formConfigs and fetches the one this matches
func findSingleFormConfig(idType, identifier, fqdn string) (*SingleFormConfig, error) {
	// TODO based on idType, we could have mtuliple forms
	sfc := &SingleFormConfig{}
	if strings.TrimSpace(identifier) == "" {
		return sfc, errors.New("identifier is empty")
	}
	sfc.FindIndex(identifier, fqdn)
	if sfc.ID == 0 {
		return sfc, errors.New("Could not find in db")
	}
	return sfc, nil
}

func newProcessedRequest(not *IncomingRequest) (*ProcessedRequest, error) {
	pr := &ProcessedRequest{
		IncomingRequest: not,
	}
	// we are ignoring scheme intentionally
	fqdn := not.Referral.Host
	formConfig, err := findSingleFormConfig(not.IDType, not.Identifier, fqdn)
	if err == nil {
		pr.SingleFormConfig = formConfig
	}
	return pr, err
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
			maxCount := min(len(ccList), ccListLimit)
			ccList = ccList[:maxCount]
			not.Cc = ccList
		} else {
			fmt.Println(err)
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
