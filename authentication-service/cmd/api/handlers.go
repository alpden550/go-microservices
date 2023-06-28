package main

import (
	"errors"
	"fmt"
	"net/http"

	helpers "github.com/alpden550/json-helpers"
)

func (app *Config) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var jsonResponse struct {
		Error   bool        `json:"error"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := helpers.ReadJSONBody(writer, request, &requestPayload); err != nil {
		err = helpers.WriteErrorJSON(writer, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		err = helpers.WriteErrorJSON(writer, errors.New("not found user"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		err = helpers.WriteErrorJSON(writer, errors.New("invalid password credentials"), http.StatusBadRequest)
		return
	}

	jsonResponse.Error = false
	jsonResponse.Message = fmt.Sprintf("Logged in user %s", user.Email)
	jsonResponse.Data = user
	err = helpers.WriteJSON(writer, http.StatusOK, jsonResponse)
	if err != nil {
		return
	}
}
