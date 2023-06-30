package main

import (
	"net/http"

	helpers "github.com/alpden550/json_helpers"
)

func (app *Config) Broker(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool

	payload := helpers.JSONResponse{
		Error:   false,
		Message: "message",
	}
	_ = tool.WriteJSON(writer, http.StatusOK, payload)
}
