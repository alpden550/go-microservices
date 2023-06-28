package main

import (
	helpers "github.com/alpden550/json-helpers"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *Config) Broker(writer http.ResponseWriter, request *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "message",
	}

	_ = helpers.WriteJSON(writer, http.StatusOK, payload)
}
