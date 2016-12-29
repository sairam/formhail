package model

const (
	// notifier end point types
	EndpointTypeEmail   = "email"
	EndpointTypeSlack   = "slack"
	EndpointTypeWebhook = "webhook"
)

// Notifier is always an outgoing notification that can be configured
type Notifier struct {
	Settings     map[string]string // any other setting like header etc.,
	EndPointURL  string            // https://.... or user@example.com
	EndPointType string            // slack,email,webhook

	// viewable fields
	Verified bool // email requires verification while slack/webhooks don't
	Internal bool // internal is an explicit one based on the registered email

	Counter
}
