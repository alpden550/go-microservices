package main

import (
	"embed"
	"fmt"
	"html/template"
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

	td := templateData{BrokerURL: app.BrokerURL}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		render(writer, "test.page.gohtml", &td)
	})

	log.Printf("Starting front end service on port %s", app.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", app.Port), nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

func render(writer http.ResponseWriter, t string, td *templateData) {
	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(writer, td); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
