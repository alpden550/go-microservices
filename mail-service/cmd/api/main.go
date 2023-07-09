package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	WebPort         string `env:"WEB_PORT"`
	MailDomain      string `env:"MAIL_DOMAIN"`
	MailHost        string `env:"MAIL_HOST"`
	MailPort        string `env:"MAIL_PORT"`
	MailUsername    string `env:"MAIL_USERNAME"`
	MailPassword    string `env:"MAIL_PASSWORD"`
	MailEncryption  string `env:"MAIL_ENCRYPTION"`
	MailFromAddress string `env:"MAIL_FROM_ADDRESS"`
	MailFromName    string `env:"MAIL_FROM_NAME"`
	Mailer          *Mail
	LoggerURL       string `env:"LOGGER_URL"`
}

func main() {
	app := Config{}
	if err := env.Parse(&app); err != nil {
		log.Fatal(err.Error())
	}

	app.Mailer = NewMail(&app)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: app.routes(),
	}

	log.Println("Started mail service at port ", app.WebPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
