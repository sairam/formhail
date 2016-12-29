package service

import (
	"net/url"

	"../model"
)

// FormRequest parses the form into schema is used an intermediary data structure
type FormRequest struct {
	Referral   *url.URL `schema:"-"`
	Identifier string   `schema:"-"`
	ReplyTo    string   `schema:"_replyto"`
	NextPage   string   `schema:"_next"`
	Subject    string   `schema:"_subject"`
	CcString   string   `schema:"_cc"`
	Format     string   `schema:"_format"`
	Gotcha     string   `schema:"_gotcha"`
}

// UserSignInService handles the business logic for signinRequest
type UserSignInService struct {
	*model.UserSignInRequest
}

// ProcessedRequest links incoming request and the config initially provided by the user
type ProcessedRequest struct {
	*model.IncomingRequest
	*model.SingleFormConfig
}

// ProcessedRequestMailContent ..
type ProcessedRequestMailContent struct {
	WebsiteURL   string
	EmailTo      string
	FormName     string
	UsersWebsite string
	UsersPage    string
	Message      *map[string][]string
}

// UserSignInRequestMail mail content
type UserSignInRequestMail struct {
	WebsiteURL  string
	EmailTo     string
	UsersDomain string
	Token       string
}

// UserSession information
type UserSession struct {
	Email  string
	Domain string
}
