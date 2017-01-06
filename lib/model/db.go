package model

// DBStore is implemented by any other store
type DBStore interface {
	load(string, string, interface{}) bool
	save(string, string, interface{}) bool
	autoincr(string) int64

	getbykey(string) string
	setbykey(string, string)
}

func getDBStore() DBStore {
	return &redisDB{}
}
