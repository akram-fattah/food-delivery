package helper

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendEmail(subject string, toEmail string, bodyTemplate string, code string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		log.Println("email credentials not set")
		return fmt.Errorf("internal error")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	body := fmt.Sprintf(bodyTemplate, code)

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{toEmail},
		[]byte(message),
	)

	if err != nil {
		log.Printf("failed to send email to %s: %v", toEmail, err)
		return fmt.Errorf("failed to send email")
	}

	return nil
}
