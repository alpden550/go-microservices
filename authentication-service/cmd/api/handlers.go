package main

import (
	"errors"
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

	jsonResponse := helpers.JSONResponse{
		Error:   false,
		Message: "fmt.Sprintf(\"Logged in user %s\", user.Email)",
		Data:    user,
	}
	err = tool.WriteJSON(writer, http.StatusOK, jsonResponse)
	if err != nil {
		return
	}
}