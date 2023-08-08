package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/getsentry/sentry-go"
)

type templateData struct {
	BrokerURL string
}

type Config struct {
	SentryDSN string `env:"SENTRY_DSN,required"`
	BrokerURL string `env:"BROKER_URL,required"`
	Port      string `env:"PORT"`
}

func main() {
	app := Config{}
	if err := env.Parse(&app); err != nil {
		log.Fatal(err.Error())
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              app.SentryDSN,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", app.Port),
		Handler: app.routes(),
	}
	log.Printf("Starting front end service on port %s", app.Port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
