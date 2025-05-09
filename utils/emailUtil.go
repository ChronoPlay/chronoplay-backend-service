package utils

import (
	"log"
	"net/smtp"
	"os"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
)

func SendEmail(to []string, subject, body string) (err *helpers.CustomEror) {
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSOWRD")

	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")

	message := []byte("Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	serr := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if serr != nil {
		return err.System(serr.Error())
	}

	log.Println("Email sent successfully")
	return nil
}
