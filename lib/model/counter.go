package model

// Counter to track no. of requests processed till ChangeTime. links to AccountLimit through AccountType
type Counter struct {
	Count      int64 // current no of requests served
	ChangeTime int64 // Next ChangeTime calculated when Count reaches the Limit
}

// GlobalCounter gets translated into Redis
type GlobalCounter struct {
	Name       string // email
	Limit      int64  // 100
	ChangeTime int64  // 1 week once counter started

	Counter
}
