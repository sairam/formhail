package main

import "fmt"

// Persist and loading of data

func (at *AccountType) load(name string) bool {
	return (&redisDB{}).load("AccountType", name, at)
}

func (at *AccountType) save() bool {
	return (&redisDB{}).save("AccountType", at.Name, at)
}

func (sfc *SingleFormConfig) load(id int) bool {
	success := (&redisDB{}).load("SingleFormConfig", fmt.Sprintf("%d", id), sfc)
	if !success {
		return false
	}
	t := &AccountType{}
	if t.load(sfc.AccountType) {
		sfc.accType = t
	}
	return success
}

func (sfc *SingleFormConfig) save() bool {
	return (&redisDB{}).save("SingleFormConfig", fmt.Sprintf("%d", sfc.ID), sfc)
}

func (sfc *SingleFormConfig) autoincr() int64 {
	return (&redisDB{}).autoincr("SingleFormConfig")
}
