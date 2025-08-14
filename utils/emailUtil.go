package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
)

func SendEmail(to []string, subject, body string) (err *helpers.CustomError) {
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSOWRD")

	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")

	// Headers + Body (HTML enabled)
	message := []byte(
		"From: " + from + "\r\n" +
			"To: " + to[0] + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
			body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	fmt.Println("Email credentials: smtpHost:", smtpHost, "from:", from, "to:", to)

	serr := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if serr != nil {
		return helpers.System(serr.Error())
	}

	log.Println("Email sent successfully")
	return nil
}
