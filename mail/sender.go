package mail

type EmailSender interface {
	SendEmail(subject string, content string, to []string, attachFiles []string) error
}
