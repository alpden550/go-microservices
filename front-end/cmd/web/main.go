package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type templateData struct {
	BrokerURL string `env:"BROKER_URL,required"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var td templateData
	if err := env.Parse(&td); err != nil {
		log.Fatal(err.Error())
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		render(writer, "test.page.gohtml", &td)
	})

	log.Println("Starting front end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(writer http.ResponseWriter, t string, td *templateData) {
	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(writer, td); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
