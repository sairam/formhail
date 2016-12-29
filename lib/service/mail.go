package service

import (
	"net/mail"

	"../common"
	"../helper"
	"../model"
	"github.com/sairam/kinli"
)

// SendEmail sends an email for a ProcessedRequest about the new notification
// TODO: plain should remove multiple line inputs by removing \n\n\n to \n
// send email for the customer once the user's website is whitelisted
func (pr *ProcessedRequest) SendEmail() {
	m := &ProcessedRequestMailContent{}
	m.Message = &pr.IncomingRequest.Message
	m.EmailTo = pr.SingleFormConfig.Email.Address
	m.WebsiteURL = common.Config.WebsiteURL
	m.FormName = pr.SingleFormConfig.Name
	m.UsersWebsite = pr.SingleFormConfig.URL
	m.UsersPage = pr.IncomingRequest.Referral.String()

	plain, _ := kinli.GetPageContent("mail_newsubmission_plain", m)
	var html string
	if helper.StringIn(pr.IncomingRequest.Format, "html") {
		html, _ = kinli.GetPageContent("mail_newsubmission", m)
	}

	e := &kinli.EmailCtx{
		From:      common.Config.FromEmail,
		To:        []*mail.Address{pr.SingleFormConfig.Email},
		Cc:        pr.IncomingRequest.Cc,
		Subject:   "New Submission for Form " + pr.SingleFormConfig.Name,
		PlainBody: plain,
		HTMLBody:  html,
	}

	e.SendEmail()
}

func (sis *UserSignInService) SendEmail() {
	sir := sis.UserSignInRequest
	m := &UserSignInRequestMail{
		WebsiteURL:  common.Config.WebsiteURL,
		EmailTo:     sir.Email,
		UsersDomain: sir.Domain,
		Token:       sir.Token,
	}
	var mailTemplate string
	if sir.RequestType == model.SirequestTypeConfirm {
		mailTemplate = "confirm"
	} else if sir.RequestType == model.SirequestTypeLogin {
		mailTemplate = "signin"
	} else {
		return
	}
	plain, _ := kinli.GetPageContent("mail_"+mailTemplate+"_plain", m)
	html, _ := kinli.GetPageContent("mail_"+mailTemplate, m)

	email, _ := mail.ParseAddress(sir.Email)

	e := &kinli.EmailCtx{
		From:      common.Config.FromEmail,
		To:        []*mail.Address{email},
		Subject:   "New SignIn Request from Domain " + sir.Domain,
		PlainBody: plain,
		HTMLBody:  html,
	}
	e.SendEmail()
}
