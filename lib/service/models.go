package service

import (
	"../model"
)

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

// UserSession information
type UserSession struct {
	Email  string
	Domain string
}
