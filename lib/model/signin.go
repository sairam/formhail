package model

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SirequestTypeConfirm = "confirmation"
	SirequestTypeLogin   = "login"
)

// UserSignInRequest is filled up when a user requests a validation/login
type UserSignInRequest struct {
	ID          int64 // auto incr key
	Email       string
	Domain      string // generic login / domain related login
	Token       string
	RequestType string // RequestType is either "confirmation" or "login"
	Status      string // used / spam / notused
	ReqTime     int64  // requested time request epoch
	ValidTime   int64  // valid time of RandomID uses time request epoch
	SEndTime    int64  // Session End Time request epoch
}

func (usir *UserSignInRequest) Load(id int) bool {
	success := (&redisDB{}).load("UserSignInRequest", fmt.Sprintf("%d", id), usir)
	return success
}

func (usir *UserSignInRequest) Save() bool {
	return (&redisDB{}).save("UserSignInRequest", fmt.Sprintf("%d", usir.ID), usir)
}

func (usir *UserSignInRequest) Autoincr() int64 {
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
func (usir *UserSignInRequest) FindIndex(token string) {
	key := strings.Join([]string{"USIRIndex", "token", token}, ":")
	t := (&redisDB{}).getKeyValue(key)
	i, err := strconv.Atoi(t)
	if err != nil || i == 0 {
		return
	}
	usir.Load(i)
}
