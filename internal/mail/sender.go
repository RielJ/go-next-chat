package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuth       = "smtp.gmail.com"
	smtpServerAddr = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachments []string,
	) error
}

type GMailSender struct {
	name     string
	from     string
	password string
}

func NewGMailSender(name, from, password string) EmailSender {
	return &GMailSender{
		name:     name,
		from:     from,
		password: password,
	}
}

func (s *GMailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachments []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.name, s.from)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, attachment := range attachments {
		_, err := e.AttachFile(attachment)
		if err != nil {
			return fmt.Errorf("failed to attach file: %w", err)
		}
	}

	smtpAuth := smtp.PlainAuth("", s.from, s.password, smtpAuth)
	err := e.Send(smtpServerAddr, smtpAuth)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
