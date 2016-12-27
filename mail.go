package main

import (
	"net/mail"
	"os"
	"strconv"

	"github.com/sairam/kinli"
)

// ProcessedRequestMailContent ..
type ProcessedRequestMailContent struct {
	WebsiteURL   string
	EmailTo      string
	FormName     string
	UsersWebsite string
	UsersPage    string
	Message      *map[string][]string
}

func init() {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil || port < 0 {
		panic("SMTP Port `" + os.Getenv("SMTP_PORT") + "` is invalid")
	}
	var smtpConfig = &kinli.EmailSMTPConfig{
		Host: os.Getenv("SMTP_HOST"),
		Port: port,
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
	}
	kinli.InitMailer(smtpConfig)
}

// send email for the first time a user attempts to send register
// func sendEmailForSigninRequest() {
// 	m := &MailContent{}
// 	plain, _ := kinli.GetPageContent("mail_signinrequest_plain", m)
// 	html, _ := kinli.GetPageContent("mail_signinrequest", m)
// 	e := &kinli.EmailCtx{
// 		From:      config.FromEmail,
// 		To:        []*mail.Address{&mail.Address{Address: "to@example.com", Name: "Mail To"}},
// 		Subject:   "HELO",
// 		PlainBody: plain,
// 		HTMLBody:  html,
// 	}
// 	e.SendEmail()
// }

// TODO: plain should remove multiple line inputs by removing \n\n\n to \n
// send email for the customer once the user's website is whitelisted
func (pr *ProcessedRequest) sendEmail() {
	m := &ProcessedRequestMailContent{}
	m.Message = &pr.IncomingRequest.Message
	m.EmailTo = pr.SingleFormConfig.Email.Address
	m.WebsiteURL = config.WebsiteURL
	m.FormName = pr.SingleFormConfig.Name
	m.UsersWebsite = pr.SingleFormConfig.URL.String()
	m.UsersPage = pr.IncomingRequest.Referral.String()

	plain, _ := kinli.GetPageContent("mail_newsubmission_plain", m)
	var html string
	if stringIn(pr.IncomingRequest.Format, "html") {
		html, _ = kinli.GetPageContent("mail_newsubmission", m)
	}

	e := &kinli.EmailCtx{
		From:      config.FromEmail,
		To:        []*mail.Address{pr.SingleFormConfig.Email},
		Cc:        pr.IncomingRequest.Cc,
		Subject:   "New Submission for Form " + pr.SingleFormConfig.Name,
		PlainBody: plain,
		HTMLBody:  html,
	}

	e.SendEmail()
}

func stringIn(list []string, key string) bool {
	for _, k := range list {
		if key == k {
			return true
		}
	}
	return false
}
