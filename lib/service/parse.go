package service

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"../model"

	"github.com/gorilla/schema"
)

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

	not := &model.IncomingRequest{}
	not.DateTime = time.Now().Unix()
	not.RemoteAddr = r.RemoteAddr

	err = r.ParseForm()
	if err != nil {
		return errors.New("Error Parsing Form")
	}
	form := new(model.FormRequest)
	decoder := schema.NewDecoder()
	decoder.Decode(form, r.PostForm)

	err = not.ParseFormFields(r.Referer(), r.RequestURI, form)
	if err != nil {
		return err
	}
	not.Message = removeRestrictedFields(r.PostForm)

	// cases - user/domain
	// 1. is not registered
	// 2. is blacklisted / limit exceeded
	// 3. has enough limit

	pr, err := newProcessedRequest(not)
	var message string
	if pr.SingleFormConfig == nil || err != nil {
		log.Println(err)
		sfc, err := newSingleFormConfig(not)
		if err != nil {
			message = "error. smth"
		}
		message = "First time user, check email"
		sendConfirmToken(sfc.Email.Address, sfc.URL)
		pr.SingleFormConfig.Confirmed = model.FormConfigUnconfirmed
		pr.Save()
	} else if pr.YetToBeConfirmed() {
		log.Println("to be confirmed")
	} else if pr.IsBlacklisted() {
		log.Println("blacklisted")
	} else if pr.DidLimitReach() {
		log.Println("Limit Reached")
	} else {
		pr.IncrIncoming()
		if !pr.Save() {
			log.Println("save notifications failed")
		} else {
			pr.SendNotifications()
		}
	}
	log.Println(message)
	return nil
}

// responsible for creating and sending an email to the user
func newSingleFormConfig(not *model.IncomingRequest) (*model.SingleFormConfig, error) {
	if not.IDType != "email" {
		return nil, errors.New("Incoming Request should be an email")
	}
	page := not.Referral.Path
	page = strings.Trim(page, "/")
	if page == "" {
		page = "homepage"
	}
	email, _ := mail.ParseAddress(not.Identifier)

	// TODO add err
	sfc := &model.SingleFormConfig{
		Name:        page,
		UID:         RandString(20),
		Email:       email,
		URL:         not.Referral.Host,
		Confirmed:   model.FormConfigRequested,
		AccountType: model.AccountTypeBasic,
		Counter: model.Counter{
			Count:      0,
			ChangeTime: 0,
		},
	}
	sfc.ID = sfc.Autoincr()
	sfc.Save()
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
		if not.EndPointType == model.EndpointTypeEmail {
			go func() { pr.SendEmail() }()
		} else if not.EndPointType == model.EndpointTypeSlack {
			pr.SendToSlack(not.EndPointURL)
		} else if not.EndPointType == model.EndpointTypeWebhook {
			pr.SendToWebhook(not.EndPointURL)
		} else {
			// not supported
		}
	}
}

// should be done with lock to avoid concurrent transactions
// looks through all formConfigs and fetches the one this matches
func findSingleFormConfig(idType, identifier, fqdn string) (*model.SingleFormConfig, error) {
	// TODO based on idType, we could have mtuliple forms
	sfc := &model.SingleFormConfig{}
	if strings.TrimSpace(identifier) == "" {
		return sfc, errors.New("identifier is empty")
	}
	sfc.FindIndex(identifier, fqdn)
	if sfc.ID == 0 {
		return sfc, errors.New("Could not find in db")
	}
	return sfc, nil
}

func newProcessedRequest(not *model.IncomingRequest) (*ProcessedRequest, error) {
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

func removeRestrictedFields(r url.Values) url.Values {
	for _, key := range []string{"_cc", "_replyto", "_next", "subject", "_format"} {
		r.Del(key)
	}
	return r
}
