package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func NewMail(app *Config) *Mail {
	port, _ := strconv.Atoi(app.MailPort)

	return &Mail{
		Domain:      app.MailDomain,
		Host:        app.MailHost,
		Port:        port,
		Username:    app.MailUsername,
		Password:    app.MailPassword,
		Encryption:  app.MailEncryption,
		FromAddress: app.MailFromAddress,
		FromName:    app.MailFromName,
	}
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}
	msg.DataMap = data

	htmlMsg, err := m.buildHTMLMessage(msg)
	if err != nil {
		log.Printf("%e", fmt.Errorf("%w", err))
		return err
	}
	textMsg, err := m.buildTextMessage(msg)
	if err != nil {
		log.Printf("%e", fmt.Errorf("%w", err))
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.
		SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, textMsg)
	email.AddAlternative(mail.TextHTML, htmlMsg)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	if err = email.Send(client); err != nil {
		return err
	}

	return nil
}

func (m *Mail) getEncryption(encryption string) mail.Encryption {
	switch encryption {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

//go:embed templates
var templateFS embed.FS

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFS(templateFS, templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMsg := tpl.String()
	formattedMsg, err = m.inlineCSS(formattedMsg)
	if err != nil {
		return "", err
	}

	return formattedMsg, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) buildTextMessage(msg Message) (string, error) {
	templateToRender := "templates/mail.text.gohtml"

	t, err := template.New("email-text").ParseFS(templateFS, templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
