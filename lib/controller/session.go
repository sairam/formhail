package controller

import (
	"os"

	"../common"
	"../service"
	"github.com/gorilla/sessions"
	"github.com/sairam/kinli"
)

func InitSession() {
	kinli.SessionStore = sessions.NewFilesystemStore("./sessions", []byte(os.Getenv("SESSION_STORE")))
	kinli.SessionName = common.Config.SessionName
	kinli.IsAuthed = isAuthed
}

func isAuthed(hc *kinli.HttpContext) bool {
	u := hc.GetSessionData("user")
	if user, ok := u.(*service.UserSession); ok && user.Email != "" {
		return true
	}
	return false
}
