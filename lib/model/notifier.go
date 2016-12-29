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

func (n *Notifier) Validate() bool {
	if n.EndPointType == EndpointTypeSlack ||
		n.EndPointType == EndpointTypeWebhook ||
		n.EndPointType == EndpointTypeEmail {
		return true
	}
	// TODO add provider specific validations like validating end point by sending an email/POST request
	return false
}

func (n *Notifier) IsInternal() bool {
	if n.EndPointType == EndpointTypeEmail {
		return true
	}
	return false
}

func (n *Notifier) RequiresVerification() bool {
	if n.EndPointType == EndpointTypeSlack || n.EndPointType == EndpointTypeWebhook {
		return false
	}
	return !n.Verified
}

func (n *Notifier) SetVerified(v bool) {
	n.Verified = v
	return
}
