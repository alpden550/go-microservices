package main

import (
	"fmt"
	"log"
	"log-service/data"
	"net/http"

	helpers "github.com/alpden550/json_helpers"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool
	var requestPayload Payload

	err := tool.ReadJSONBody(writer, request, &requestPayload)
	if err != nil {
		log.Printf("%e", fmt.Errorf("%w", err))
		return
	}

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	if err = app.Models.LogEntry.InsertRecord(event); err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	response := helpers.JSONResponse{
		Error:   false,
		Message: "logged",
	}

	if err = tool.WriteJSON(writer, http.StatusAccepted, response); err != nil {
		log.Printf("%e", fmt.Errorf("%w", err))
		return
	}
}
