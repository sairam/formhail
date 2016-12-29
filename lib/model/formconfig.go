package model

import (
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

const (
	FormConfigRequested   = "requested"   // someone has requested this, we are yet to send the email
	FormConfigUnconfirmed = "unconfirmed" // unconfirmed means that we have sent the email, but its not yet confirmed
	FormConfigConfirmed   = "confirmed"   // confirmed means the email was clicked
	FormConfigSpam        = "spam"        // unconfirmed can transition to spam/confirmed

	IDTypeEmail = "email"
	IDTypeUID   = "uid"
)

// SingleFormConfig has configuration for a single form
type SingleFormConfig struct {
	ID            int64 // unique internal identifier, auto incrementer
	Name          string
	UID           string // UID is an alias to the Email, but a randomly generated string
	Email         *mail.Address
	URL           string // A page means a single page is supported
	URLType       string // URLType can be page or domain or regexp for URL that needs to be matched
	Confirmed     string // formConfigRequested, formConfigConfirmed, formConfigUnconfirmed, formConfigSpam
	ConfirmedDate string // datetime at which this confirmation was made

	// Counters to track for incoming
	Counter // TODO incoming counter should be at Domain or Email level instead of form level

	AccountType string
	accType     *AccountType // Links to an Account Type via the string

	// All notifications to external points can be configured through this
	// Limits apply based on AccountType
	Notifications map[string]*Notifier // default outgoing notification is added on confirmation
}

func (sfc *SingleFormConfig) Load(id int) bool {
	success := (&redisDB{}).load("SingleFormConfig", fmt.Sprintf("%d", id), sfc)
	if !success {
		return false
	}
	t := &AccountType{}
	if t.Load(sfc.AccountType) {
		sfc.accType = t
	}
	return success
}

func (sfc *SingleFormConfig) Save() bool {
	return (&redisDB{}).save("SingleFormConfig", fmt.Sprintf("%d", sfc.ID), sfc)
}

func (sfc *SingleFormConfig) Autoincr() int64 {
	return (&redisDB{}).autoincr("SingleFormConfig")
}

// Index saves the email/domain mapping
func (sfc *SingleFormConfig) Index() {
	email := sfc.Email.Address
	domain := sfc.URL

	key := strings.Join([]string{"SFCIndex", "domemail", email, domain}, ":")
	id := fmt.Sprintf("%d", sfc.ID)
	(&redisDB{}).setKeyValue(key, id)
}

// FindIndex locates and loads the config
func (sfc *SingleFormConfig) FindIndex(email, domain string) {
	key := strings.Join([]string{"SFCIndex", "domemail", email, domain}, ":")
	t := (&redisDB{}).getKeyValue(key)
	i, err := strconv.Atoi(t)
	if err != nil || i == 0 {
		return
	}
	sfc.Load(i)
}

// YetToBeConfirmed ..
func (c *SingleFormConfig) YetToBeConfirmed() bool {
	return c.Confirmed == FormConfigUnconfirmed || c.Confirmed == FormConfigRequested
}

// IsBlacklisted ..
func (c *SingleFormConfig) IsBlacklisted() bool {
	return c.Confirmed == FormConfigSpam
}

// DidLimitReach checks if we reached the limit for the account?
// checks for incoming requests
// TODO verify the limit based on the account type
// incr with a lock and save to store
func (c *SingleFormConfig) DidLimitReach() bool {
	// no need to change the time
	if c.ChangeTime == 0 {
		c.ChangeTime = time.Now().Unix() - 10
	}
	accLimit := c.accType.Limits["incoming:form"]
	currentTime := time.Now().Unix()
	for currentTime > c.ChangeTime {
		c.ChangeTime += accLimit.Period
		c.Count = 0
	}
	if c.Count < accLimit.Limit {
		return false
	}
	return true
}

// IncrIncoming usage field
// TODO incr with a lock
func (c *SingleFormConfig) IncrIncoming() {
	c.Count = c.Count + 1
}
