package main

import (
	"fmt"
	"log"
	"net/smtp"
)

// https://github.com/spf13/viper

// also see https://gist.github.com/chrisgillis/10888032

func main() {

	// 1. Listen on a port
	// 2. Get all the form details / parse the form
	// 3. send an email with the form name and the key/value attributes

}

/*
MVP Next Step:
1. Register for a user via a random generated id / string
2. link the user's email address to send to
3. User's website to match
4. Limit on no. of messages per hour (account limit) and per IP (spam limit)
5. Send emails to users.
6. Support Attachments with limit.
7. If attachment size > 10MB, put in s3 and provide unrevokable link via s3 signing
*/

type plainSMTP struct{}

const fromEmail = "from@example.com"
const toEmail = "to@example.com"
const connectionHost = "localhost:25"

func (*plainSMTP) Send() {

	body = "This is the email body"
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(connectionHost)
	if err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient first
	if err = c.Mail(fromEmail); err != nil {
		log.Fatal(err)
	}
	if err = c.Rcpt(toEmail); err != nil {
		log.Fatal(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, body)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
