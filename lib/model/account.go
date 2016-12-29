package model

const AccountTypeBasic = "basic"

// AccountType has a name, description and limits based on the type of channel
type AccountType struct {
	Name        string // Basic
	Description string
	Limits      map[string]AccountLimit // Has different Configuration
}

// AccountLimit defines how many requests can be accepted per a period
type AccountLimit struct {
	Type string // incoming, outgoing:slack, outgoing:email, outgoing:webhook
	// Limit & Period are configurable at Account / User level
	// if limit is -1, unlimited will be sent.
	Limit  int64 // no. of Requests to limit to until ChangeTime
	Period int64 // no. of seconds from ChangeTime it will reset to ChangeTime += Period & Count = 0
}

func (at *AccountType) Load(name string) bool {
	return (&redisDB{}).load("AccountType", name, at)
}

func (at *AccountType) Save() bool {
	return (&redisDB{}).save("AccountType", at.Name, at)
}
