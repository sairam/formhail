package main

import (
	"fmt"
	"strconv"
	"strings"
)

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

// Index saves the email/domain mapping
func (sfc *SingleFormConfig) Index() {
	email := sfc.Email.Address
	domain := sfc.URL.String()

	key := strings.Join([]string{"SFCIndex", "domemail", email, domain}, ":")
	id := fmt.Sprintf("%d", sfc.ID)
	(&redisDB{}).setKeyValue(key, id)
}

// FindIndex locates and loads the config
func (sfc *SingleFormConfig) FindIndex(email, domain string) {
	key := strings.Join([]string{"SFCIndex", "domemail", email, domain}, ":")
	t := (&redisDB{}).getKeyValue(key)
	i, err := strconv.Atoi(t)
	fmt.Println(i)
	fmt.Println(err)
	if err != nil || i == 0 {
		return
	}
	sfc.load(i)
}

// UserSignInRequest persist functions
func (usir *UserSignInRequest) load(id int) bool {
	success := (&redisDB{}).load("UserSignInRequest", fmt.Sprintf("%d", id), usir)
	return success
}

func (usir *UserSignInRequest) save() bool {
	return (&redisDB{}).save("UserSignInRequest", fmt.Sprintf("%d", usir.ID), usir)
}

func (usir *UserSignInRequest) autoincr() int64 {
	return (&redisDB{}).autoincr("UserSignInRequest")
}

// Index indexes data
func (usir *UserSignInRequest) Index() {
	token := usir.Token

	key := strings.Join([]string{"USIRIndex", "token", token}, ":")
	id := fmt.Sprintf("%d", usir.ID)
	(&redisDB{}).setKeyValue(key, id)
	// TODO: set expiry based on the request data based usir.ValidTime
}

// FindIndex finds based on the indexed data
func (usir UserSignInRequest) FindIndex(token string) {
	key := strings.Join([]string{"USIRIndex", "token", token}, ":")
	t := (&redisDB{}).getKeyValue(key)
	i, err := strconv.Atoi(t)
	if err != nil || i == 0 {
		return
	}
	usir.load(i)

}
