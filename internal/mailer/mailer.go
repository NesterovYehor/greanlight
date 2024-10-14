package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

// Embed template files from the "templates" directory.
//
//go:embed templates/*
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(port int, host, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = time.Second * 5

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (mailer *Mailer) Send(recipient, templateFile string, data any) error {
	// Parse the templates from the embedded filesystem.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Prepare the subject, plain text body, and HTML body.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Create a new email message.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", mailer.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// Send the email.
	err = mailer.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}

