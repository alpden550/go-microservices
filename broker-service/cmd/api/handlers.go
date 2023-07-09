package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	helpers "github.com/alpden550/json_helpers"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool

	payload := helpers.JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = tool.WriteJSON(writer, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool
	var requestPayload RequestPayload

	if err := tool.ReadJSONBody(writer, request, &requestPayload); err != nil {
		err = tool.WriteErrorJSON(writer, err)
		if err != nil {
			return
		}
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(writer, requestPayload.Auth)
	case "log":
		app.logItem(writer, requestPayload.Log)
	case "mail":
		app.sendMail(writer, requestPayload.Mail)
	default:
		err := tool.WriteErrorJSON(writer, errors.New("unknown action"))
		if err != nil {
			return
		}
	}
}

func (app *Config) authenticate(writer http.ResponseWriter, a AuthPayload) {
	var tool helpers.Tool
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	url := fmt.Sprintf("%s/authenticate", app.AuthURL)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}
	defer response.Body.Close()

	var jsonFromService helpers.JSONResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	if jsonFromService.Error {
		_ = tool.WriteErrorJSON(writer, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	payload := helpers.JSONResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}

	err = tool.WriteJSON(writer, http.StatusAccepted, payload)
	if err != nil {
		return
	}
}

func (app *Config) logItem(writer http.ResponseWriter, entry LogPayload) {
	var tool helpers.Tool
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	loggerURL := fmt.Sprintf("%s/log", app.LoggerURL)
	request, err := http.NewRequest("POST", loggerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	payload := helpers.JSONResponse{
		Error:   false,
		Message: "logged",
	}
	err = tool.WriteJSON(writer, http.StatusAccepted, payload)
	if err != nil {
		log.Printf("%e", fmt.Errorf("%e", err))
		return
	}
}

func (app *Config) sendMail(writer http.ResponseWriter, msg MailPayload) {
	var tool helpers.Tool
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailerURL := fmt.Sprintf("%s/send", app.MailerURL)
	request, err := http.NewRequest("POST", mailerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		_ = tool.WriteErrorJSON(writer, errors.New("error calling mail service"))
		return
	}

	payload := helpers.JSONResponse{
		Error:   false,
		Message: "Message sent to " + msg.To,
	}

	err = tool.WriteJSON(writer, http.StatusAccepted, payload)
	if err != nil {
		return
	}
}
