package main

import (
	"net/http"

	helpers "github.com/alpden550/json-helpers"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *Config) Broker(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool

	payload := jsonResponse{
		Error:   false,
		Message: "message",
	}

	_ = tool.WriteJSON(writer, http.StatusOK, payload)
}
