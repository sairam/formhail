package main

import (
	"errors"
	"net/http"
	"net/mail"
	"net/url"
	"time"

	"github.com/sairam/kinli"
)

// if type=Confirm in UserSignInRequest, use the domain
func verifyUserConfirmToken(token string) (user *UserSession, err error) {
	sir := &UserSignInRequest{}
	sir.FindIndex(token)
	if sir.ID == 0 || sir.RequestType != sirequestTypeConfirm {
		return nil, errors.New("Could not find request")
	}
	// find request and assign
	sfc := &SingleFormConfig{}
	sfc.FindIndex(sir.Email, sir.Domain)
	if sfc.Confirmed == formConfigConfirmed {
		return nil, errors.New("User/Token is already confirmed")
	}
	sfc.Confirmed = formConfigConfirmed
	sfc.ConfirmedDate = time.Now().String()
	if sfc.Notifications == nil {
		sfc.Notifications = make(map[string]*Notifier)
	}
	sfc.Notifications["email"] = &Notifier{
		EndPointURL:  sfc.Email.Address,
		EndPointType: "email",
		Verified:     true,
		Internal:     true,
	}
	sfc.save()

	user = &UserSession{
		Email:  sir.Email,
		Domain: sir.Domain,
	}
	return user, nil
}

func verifyLoginToken(token string) (user *UserSession, err error) {
	sir := &UserSignInRequest{}
	sir.FindIndex(token)
	if sir.ID == 0 || sir.RequestType != sirequestTypeLogin {
		return nil, errors.New("Could not find request")
	}
	user = &UserSession{
		Email:  sir.Email,
		Domain: sir.Domain,
	}
	return user, nil
}

func sendConfirmToken(email, domain string) {
	userSignInRequest, _ := makeUserSignInRequest(email, domain, sirequestTypeConfirm)
	userSignInRequest.sendEmail()
}

// makeAToken is the POST call that will send an email
func makeAToken(r *http.Request) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}

	email, err := mail.ParseAddress(r.Form["email"][0])
	if err != nil {
		return
	}

	parsed, err := url.Parse(r.Form["domain"][0])
	domain := parsed.Host

	userSignInRequest, err := makeUserSignInRequest(email.Address, domain, sirequestTypeLogin)
	userSignInRequest.sendEmail()

	return
}

// UserSignInRequestMail mail content
type UserSignInRequestMail struct {
	WebsiteURL   string
	EmailTo      string
	UsersDomain  string
	Confirmation string
}

func (sir *UserSignInRequest) sendEmail() {
	m := &UserSignInRequestMail{
		WebsiteURL:   config.WebsiteURL,
		EmailTo:      sir.Email,
		UsersDomain:  sir.Domain,
		Confirmation: sir.Token,
	}
	var mailTemplate string
	if sir.RequestType == sirequestTypeConfirm {
		mailTemplate = "confirm"
	} else if sir.RequestType == sirequestTypeLogin {
		mailTemplate = "signin"
	} else {
		return
	}
	plain, _ := kinli.GetPageContent("mail_"+mailTemplate+"_plain", m)
	html, _ := kinli.GetPageContent("mail_"+mailTemplate, m)

	email, _ := mail.ParseAddress(sir.Email)

	e := &kinli.EmailCtx{
		From:      config.FromEmail,
		To:        []*mail.Address{email},
		Subject:   "New SignIn Request from Domain " + sir.Domain,
		PlainBody: plain,
		HTMLBody:  html,
	}
	e.SendEmail()
}

func makeUserSignInRequest(email, domain, requestType string) (*UserSignInRequest, error) {
	currTime := time.Now().Unix()
	oneDay := int64(60 * 60 * 24)
	usir := &UserSignInRequest{
		Email:       email,
		Domain:      domain,
		RequestType: requestType,
		Token:       RandString(40),
		Status:      "notused",
		ReqTime:     currTime,
		ValidTime:   currTime + oneDay, // 1 day
		SEndTime:    oneDay * 30,       // 30 days
	}
	usir.ID = usir.autoincr()
	return usir, nil
}
