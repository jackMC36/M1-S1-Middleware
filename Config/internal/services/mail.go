package services

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/sirupsen/logrus"
)

// SendMail sends an email to `to` with `subject` and `body`.
// If SMTP environment variables are not set, the function will log the message instead.
func SendMail(to, subject, body string) error {
	host := os.Getenv("MAIL_SMTP_HOST")
	port := os.Getenv("MAIL_SMTP_PORT")
	user := os.Getenv("MAIL_SMTP_USER")
	pass := os.Getenv("MAIL_SMTP_PASS")

	if host == "" || port == "" {
		// Fallback: just log
		logrus.Infof("[MAIL-DRYRUN] To: %s | Subject: %s | Body: %s", to, subject, body)
		return nil
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", user, pass, host)

	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/plain; charset=UTF-8;\r\n\r\n" +
		body

	if err := smtp.SendMail(addr, auth, user, []string{to}, []byte(msg)); err != nil {
		return err
	}
	logrus.Infof("Email sent to %s", to)
	return nil
}
