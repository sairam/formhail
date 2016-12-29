package formtoemail

import (
	"os"
	"strconv"

	"github.com/sairam/kinli"
)

// InitMail starts the email service
func InitMail() {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil || port < 0 {
		panic("SMTP Port `" + os.Getenv("SMTP_PORT") + "` is invalid")
	}
	var smtpConfig = &kinli.EmailSMTPConfig{
		Host: os.Getenv("SMTP_HOST"),
		Port: port,
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
	}
	kinli.InitMailer(smtpConfig)
}
