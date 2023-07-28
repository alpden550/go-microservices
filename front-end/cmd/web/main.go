package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
)

type templateData struct {
	BrokerURL string `env:"BROKER_URL,required"`
}

func main() {
	var td templateData
	if err := env.Parse(&td); err != nil {
		log.Fatal(err.Error())
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		render(writer, "test.page.gohtml", &td)
	})

	log.Println("Starting front end service on port 8000")
	err := http.ListenAndServe(":8000", nil)
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
