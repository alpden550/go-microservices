package main

import (
	"errors"
	"fmt"
	helpers "github.com/alpden550/json-helpers"
	"net/http"
)

func (app *Config) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool

	var jsonResponse struct {
		Error   bool        `json:"error"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := tool.ReadJSONBody(writer, request, &requestPayload); err != nil {
		err = tool.WriteErrorJSON(writer, err)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		err = tool.WriteErrorJSON(writer, errors.New("not found user"))
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		err = tool.WriteErrorJSON(writer, errors.New("invalid password credentials"))
		return
	}

	jsonResponse.Error = false
	jsonResponse.Message = fmt.Sprintf("Logged in user %s", user.Email)
	jsonResponse.Data = user
	err = tool.WriteJSON(writer, http.StatusOK, jsonResponse)
	if err != nil {
		return
	}
}
