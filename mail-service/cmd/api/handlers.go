package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	helpers "github.com/alpden550/json_helpers"
)

func (app *Config) SendMail(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	if err := tool.ReadJSONBody(writer, request, &requestPayload); err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err := app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Print("PANIC", err)
		return
	}

	err = app.logRequest(fmt.Sprintf("Sent email to %s", msg.To), fmt.Sprintf("%s", msg.Data))
	if err != nil {
		log.Printf("%e", fmt.Errorf("%e", err))
	}

	response := helpers.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("sent email message to %s", requestPayload.To),
	}
	err = tool.WriteJSON(writer, http.StatusAccepted, response)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := fmt.Sprintf("%s/log", app.LoggerURL)
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
