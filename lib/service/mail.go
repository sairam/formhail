package service

import (
	"net/mail"

	"../common"
	"../helper"
	"github.com/sairam/kinli"
)

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
