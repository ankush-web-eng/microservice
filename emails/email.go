package email

import (
	"fmt"

	"github.com/ankush-web-eng/microservice/config"
	"gopkg.in/gomail.v2"
)

type EmailDetails struct {
	From    string
	To      []string
	Subject string
	Body    string
}

type EmailDetailsAsService struct {
	From     string
	To       []string
	Subject  string
	Body     string
	Username string
	Password string
}

func SendEmail(details EmailDetails) error {
	smtpConfig := config.LoadSMTPConfig()

	m := gomail.NewMessage()
	m.SetHeader("From", details.From)
	m.SetHeader("To", details.To...)
	m.SetHeader("Subject", details.Subject)
	m.SetBody("text/plain", details.Body)

	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}
	return nil
}

func SendEmailAsService(details EmailDetailsAsService) error {
	smtpConfig := config.LoadSMTPConfigAsService(config.SMTPConfigAsService{Username: details.Username, Password: details.Password})

	m := gomail.NewMessage()
	m.SetHeader("From", smtpConfig.Username)
	m.SetHeader("To", details.To...)
	m.SetHeader("Subject", details.Subject)
	m.SetBody("text/plain", details.Body)

	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}
	return nil
}
