package service

import (
	"errors"
	"net/http"
	"net/mail"
	"net/url"
	"time"

	"../helper"
	"../model"
)

const (
	oneDay = int64(60 * 60 * 24)
)

// if type=Confirm in UserSignInRequest, use the domain
func verifyUserConfirmToken(token string) (user *UserSession, err error) {
	sir := &model.UserSignInRequest{}
	sir.FindIndex(token)
	if sir.ID == 0 || sir.RequestType != model.SirequestTypeConfirm {
		return nil, errors.New("Could not find request")
	}

	// find request and assign
	sfc := &model.SingleFormConfig{}
	sfc.FindIndex(sir.Email, sir.Domain)
	if sfc.Confirmed == model.FormConfigConfirmed {
		return nil, errors.New("User/Token is already confirmed")
	}

	sfc.Confirmed = model.FormConfigConfirmed
	sfc.ConfirmedDate = time.Now().String()

	if sfc.Notifications == nil {
		sfc.Notifications = make(map[string]*model.Notifier)
	}
	sfc.Notifications["email"] = &model.Notifier{
		EndPointURL:  sfc.Email.Address,
		EndPointType: "email",
		Verified:     true,
		Internal:     true,
	}
	sfc.Save()

	user = &UserSession{
		Email:  sir.Email,
		Domain: sir.Domain,
	}
	return user, nil
}

func verifyLoginToken(token string) (user *UserSession, err error) {
	sir := &model.UserSignInRequest{}
	sir.FindIndex(token)
	if sir.ID == 0 || sir.RequestType != model.SirequestTypeLogin {
		return nil, errors.New("Could not find request")
	}
	user = &UserSession{
		Email:  sir.Email,
		Domain: sir.Domain,
	}
	return user, nil
}

func sendConfirmToken(email, domain string) {
	userSignInRequest, _ := newUserSignInRequest(email, domain, model.SirequestTypeConfirm)
	userSignInRequest.Save()
	userSignInRequest.Index()
	(&UserSignInService{UserSignInRequest: userSignInRequest}).SendEmail()
}

// requestToAuthenticate takes in the email/domain and possibly captcha
// and makes a token which the user can enter or click to login
// This is the password less mechanism that takes care of on-demand authentication
func requestToAuthenticate(r *http.Request) (err error) {
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

	userSignInRequest, err := newUserSignInRequest(email.Address, domain, model.SirequestTypeLogin)
	userSignInRequest.Save()
	userSignInRequest.Index()
	(&UserSignInService{UserSignInRequest: userSignInRequest}).SendEmail()

	return
}

func newUserSignInRequest(email, domain, requestType string) (*model.UserSignInRequest, error) {
	currTime := time.Now().Unix()
	usir := &model.UserSignInRequest{
		Email:       email,
		Domain:      domain,
		RequestType: requestType,
		Token:       helper.RandString(40),
		Status:      "notused",
		ReqTime:     currTime,
		ValidTime:   currTime + oneDay, // 1 day
		SEndTime:    oneDay * 30,       // 30 days
	}
	usir.ID = usir.Autoincr()
	return usir, nil
}
