package utils

import (
	"fmt"
	"net/smtp"
)

func SendResetEmail(toEmail, resetLink string) error {
	from := "no-reply@yourapp.com"
	password := "APP_EMAIL_PASSWORD"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := fmt.Sprintf(`From: Your App <%s>
To: %s
Subject: Reset Password

Click the link below to reset your password:

%s

This link will expire in 15 minutes.
`, from, toEmail, resetLink)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{toEmail},
		[]byte(message),
	)
}
