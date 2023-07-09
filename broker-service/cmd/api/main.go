package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
)

const port = "80"

type Config struct {
	AuthURL   string `env:"AUTH_URL,required"`
	LoggerURL string `env:"LOGGER_URL,required"`
	MailerURL string `env:"MAILER_URL,required"`
}

func main() {
	app := Config{}
	if err := env.Parse(&app); err != nil {
		log.Fatal(err.Error())
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
