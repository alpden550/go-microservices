package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

func (app *Config) IndexPage(writer http.ResponseWriter, request *http.Request) {
	td := templateData{BrokerURL: app.BrokerURL}
	render(writer, "test.page.gohtml", &td)
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
