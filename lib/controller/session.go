package controller

import (
	"os"

	"../common"
	"../service"
	"github.com/gorilla/sessions"
	"github.com/sairam/kinli"
)

// InitSession starts the configures the session handler
func InitSession() {
	kinli.SessionStore = sessions.NewFilesystemStore("./sessions", []byte(os.Getenv("SESSION_STORE")))
	kinli.SessionName = common.Config.SessionName
	kinli.IsAuthed = service.IsUserAuthed
}
