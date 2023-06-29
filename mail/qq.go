package mail

import (
	"gopkg.in/gomail.v2"
)

var _ EmailSender = &QQSender{}

const (
	smtpAuthAddressWithQQ = "smtp.qq.com"
	smtpServerPortWithQQ  = 465
)

type QQSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewQQSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &QQSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *QQSender) SendEmail(subject string, content string, to []string, attachFiles []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", sender.fromEmailAddress)
	m.SetHeader("To", to...)
	m.SetAddressHeader("Cc", sender.fromEmailAddress, sender.name)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)
	for _, file := range attachFiles {
		m.Attach(file)
	}

	d := gomail.NewDialer(smtpAuthAddressWithQQ, smtpServerPortWithQQ, sender.fromEmailAddress, sender.fromEmailPassword)

	// Send the email to Bob, Cora and Dan.
	return d.DialAndSend(m)
}
