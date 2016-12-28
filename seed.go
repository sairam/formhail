package main

type seed struct{}

// seed AccountType
// (&seed{}).accountTypes()
func (*seed) accountTypes() {
	var at = &AccountType{}
	at.Name = "basic" // plan name
	var als = []AccountLimit{
		AccountLimit{
			Type:   "incoming:form",
			Limit:  200,
			Period: 86400 * 7, // seconds
		},
		AccountLimit{
			Type:   "outgoing:email",
			Limit:  200,
			Period: 86400 * 7, // seconds
		},
		AccountLimit{
			Type:   "outgoing:slack",
			Limit:  100000,
			Period: 86400 * 7, // seconds
		},
	}
	at.Limits = make(map[string]AccountLimit)
	for _, al := range als {
		at.Limits[al.Type] = al
	}
	at.save()
}
