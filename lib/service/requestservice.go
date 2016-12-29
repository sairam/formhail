package service

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"../helper"
	"../model"

	"github.com/gorilla/schema"
)

const (
	ccListLimit = 4
)

var (
	restrictedFieldList = []string{"_cc", "_replyto", "_next", "subject", "_format"}

	ErrParsingForm              = errors.New("Form could not be parsed")
	ErrFormIdentifierNotPresent = errors.New("identifier is empty")
	ErrFormConfigNotPresent     = errors.New("Form Config not present")
	ErrShouldBeEmail            = errors.New("Incoming Request should be an email")
	ErrReferralShouldBePresent  = errors.New("Referral is blank")
	ErrItsABot                  = errors.New("Gotcha is set. Expected not to exist")
	ErrReqURINotParsable        = errors.New("Request URI is not parsable")
)

// FormSubmissionRequest handles the main traffic of consuming requests
func FormSubmissionRequest(w http.ResponseWriter, r *http.Request) {
	err := formSubmissionParser(w, r)
	if err != nil {
		http.NotFound(w, r)
		// TODO: put in a better error message
		return
	}
}

// TODO: any request of > 10KB size should be rejected
func formSubmissionParser(w http.ResponseWriter, r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return ErrParsingForm
	}

	form := new(model.FormRequest)
	decoder := schema.NewDecoder()
	decoder.Decode(form, r.PostForm)

	not, err := parseFormFields(r.Referer(), r.RequestURI, form)
	if err != nil {
		return err
	}

	not.DateTime = time.Now().Unix()
	not.RemoteAddr = r.RemoteAddr
	not.Message = removeRestrictedFields(r.PostForm)
	return formSubmissionLogic(not)
}

func formSubmissionLogic(not *model.IncomingRequest) error {
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
		// TODO show a page which says awaiting confirmation and a form submission
		// to request again with a captcha
	} else if pr.IsBlacklisted() {
		log.Println("blacklisted")
		// TODO - show the relevant page
	} else if pr.DidLimitReach() {
		log.Println("Limit Reached")
		// TODO - show a page or fail silently
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
		return nil, ErrShouldBeEmail
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
		UID:         helper.RandString(20),
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
	// TODO validate *identifier* . it should be a valid email address or alpha since those are two types we generate
	// else we may be susceptible to redis injections since we delimit by ":"
	sfc := &model.SingleFormConfig{}
	if strings.TrimSpace(identifier) == "" {
		return sfc, ErrFormIdentifierNotPresent
	}
	sfc.FindIndex(identifier, fqdn)
	if sfc.ID == 0 {
		return sfc, ErrFormConfigNotPresent
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
	for _, key := range restrictedFieldList {
		r.Del(key)
	}
	return r
}

// ParseFormFields parses the input data and fills in the struct
func parseFormFields(referrer, requestURI string, form *model.FormRequest) (*model.IncomingRequest, error) {
	var not = &model.IncomingRequest{}

	var err error

	form.Referral, err = url.Parse(referrer)
	if err != nil {
		return not, ErrReferralShouldBePresent
	}
	not.Referral = form.Referral

	uri, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return not, ErrReqURINotParsable
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
		return not, ErrItsABot
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

	return not, nil
}
