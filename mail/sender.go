// simple-bank/mail/sender.go
package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	gmailSMTPServer = "smtp.gmail.com"
	gmailSMTPPort   = 587
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(
	name string,
	fromEmailAddress string,
	fromEmailPassword string,
) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, file := range attachFiles {
		if _, err := e.AttachFile(file); err != nil {
			return err
		}
	}

	auth := smtp.PlainAuth(
		"",
		sender.fromEmailAddress,
		sender.fromEmailPassword,
		gmailSMTPServer,
	)

	addr := fmt.Sprintf("%s:%d", gmailSMTPServer, gmailSMTPPort)

	return e.Send(addr, auth)
}
