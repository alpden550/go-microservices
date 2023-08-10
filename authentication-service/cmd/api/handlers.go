package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	helpers "github.com/alpden550/json_helpers"
)

func (app *Config) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := tool.ReadJSONBody(writer, request, &requestPayload); err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, errors.New("not found user"))
		return
	}

	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		_ = tool.WriteErrorJSON(writer, errors.New("invalid password credentials"))
		return
	}

	err = app.logRequest("Authenticated", fmt.Sprintf("logged in %s", requestPayload.Email))
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	jsonResponse := helpers.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	err = tool.WriteJSON(writer, http.StatusAccepted, jsonResponse)
	if err != nil {
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
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Client.Do(request)
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted || err != nil {
		return err
	}

	return nil
}
