package mail

import (
	"gopkg.in/gomail.v2"
)

var _ EmailSender = &GmailSender{}

const (
	smtpAuthAddressWithGmail = "smtp.gmail.com"
	smtpServerPortWithGmail  = 587
)

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func (sender *GmailSender) SendEmail(subject string, content string, to []string, attachFiles []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", sender.fromEmailAddress)
	m.SetHeader("To", to...)
	m.SetAddressHeader("Cc", sender.fromEmailAddress, sender.name)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)
	for _, file := range attachFiles {
		m.Attach(file)
	}

	d := gomail.NewDialer(smtpAuthAddressWithGmail, smtpServerPortWithGmail, sender.fromEmailAddress, sender.fromEmailPassword)

	// Send the email to Bob, Cora and Dan.
	return d.DialAndSend(m)
}
func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}
