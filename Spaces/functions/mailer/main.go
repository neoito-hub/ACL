package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
)

// Get the smtp auth data
func getSmtpAuth() (smtp.Auth, string, string) {
	// Sender data.

	from := os.Getenv("SHIELD_MAILER_EMAIL")
	password := os.Getenv("SHIELD_MAILER_PASSWORD")

	// smtp server configuration.
	smtpHost := os.Getenv("SHIELD_MAILER_HOST")
	smtpPort := os.Getenv("SHIELD_MAILER_PORT")

	smtpUrl := smtpHost + ":" + smtpPort

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	return auth, smtpUrl, from
}

// Send  email verification and return success/error message
func SendEmail(body bytes.Buffer, receiverEmail []string) (string, error) {

	// Get mail template body
	to := receiverEmail

	// Get smtp auth data
	auth, smtpUrl, from := getSmtpAuth()

	fmt.Println(auth, smtpUrl, from, to)

	// Sending email.
	err := smtp.SendMail(smtpUrl, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return "Error", err
	}

	fmt.Println("Email Sent!")
	return "Done", err
}
