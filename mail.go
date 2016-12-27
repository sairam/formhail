package main

import (
	"net/mail"
	"os"
	"strconv"

	"github.com/sairam/kinli"
)

// MailContent ..
type MailContent struct {
	WebsiteURL string
	Content    interface{}
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
func sendEmailForSigninRequest() {
	m := &MailContent{}
	plain, _ := kinli.GetPageContent("mail_signinrequest_plain", m)
	html, _ := kinli.GetPageContent("mail_signinrequest", m)
	e := &kinli.EmailCtx{
		From:      &mail.Address{Address: "from@example.com", Name: "Mailer"},
		To:        []*mail.Address{&mail.Address{Address: "to@example.com", Name: "Mail To"}},
		Subject:   "HELO",
		PlainBody: plain,
		HTMLBody:  html,
	}
	e.SendEmail()
}

// send email for the customer once the user's website is whitelisted
func sendEmailToCustomer() {
}
