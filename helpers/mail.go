package helpers

import (
	"net/smtp"
	"os"
)

func SendMail(to, subject, content string) error {
	// Sender data.
	from := os.Getenv("Email")
	password := os.Getenv("EmailPassword")

	// Receiver email address.
	tos := []string{to}

	// Message
	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"" +
		"\r\n" +
		content + "\r\n"

	// smtp server configuration.
	host := "mail.enesonmez.com"
	port := "587"
	// Authentication.
	auth := smtp.PlainAuth("", from, password, host)
	// Sending email.
	address := host + ":" + port
	err := smtp.SendMail(address, auth, from, tos, []byte(msg))
	return err
}
