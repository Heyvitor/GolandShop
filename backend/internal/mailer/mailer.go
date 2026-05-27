package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type Mailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func New(host string, port int, username, password, from string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *Mailer) SendWelcomeEmail(toEmail, toName string) error {
	subject := "Bem-vindo ao nosso sistema!"
	body := fmt.Sprintf("Olá %s,\n\nSeja muito bem-vindo ao nosso sistema! Estamos felizes em ter você aqui.\n\nAbraços,\nEquipe GoApi", toName)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", m.from, toEmail, subject, body)

	// If using something like mailhog/mailtrap locally, auth might be optional, but we pass it anyway
	var auth smtp.Auth
	if m.username != "" && m.password != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	return smtp.SendMail(addr, auth, m.from, []string{toEmail}, []byte(msg))
}

func (m *Mailer) SendHTMLTemplate(toEmail, subject string, tmpl *template.Template, data any) error {
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	header := make(map[string]string)
	header["From"] = m.from
	header["To"] = toEmail
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	var msg bytes.Buffer
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.Write(body.Bytes())

	var auth smtp.Auth
	if m.username != "" && m.password != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	return smtp.SendMail(addr, auth, m.from, []string{toEmail}, msg.Bytes())
}
