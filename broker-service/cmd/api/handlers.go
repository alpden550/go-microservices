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
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	log.Printf("%#v", jsonFromService)

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
