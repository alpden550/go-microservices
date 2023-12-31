package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"broker/event"
	"broker/logs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

type RPCPayload struct {
	Name string
	Data string
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
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(writer, requestPayload.Auth)
	case "log":
		app.logEventRabbit(writer, requestPayload.Log)
	case "log-rpc":
		app.logEventViaRPC(writer, requestPayload.Log)
	case "mail":
		app.sendMail(writer, requestPayload.Mail)
	default:
		_ = tool.WriteErrorJSON(writer, errors.New("unknown action"))
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

func (app *Config) logEventRabbit(writer http.ResponseWriter, l LogPayload) {
	tool := helpers.Tool{}
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	payload := helpers.JSONResponse{
		Error:   false,
		Message: "logged via rabbitMQ",
	}
	err = tool.WriteJSON(writer, http.StatusAccepted, payload)
	if err != nil {
		return
	}

}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	jsonPayload, _ := json.MarshalIndent(&payload, "", "\t")
	if err = emitter.Push(string(jsonPayload), "log.INFO"); err != nil {
		return err
	}

	return nil
}

func (app *Config) logEventViaRPC(writer http.ResponseWriter, l LogPayload) {
	tool := helpers.Tool{}

	client, err := rpc.Dial("tcp", fmt.Sprintf("%s", app.RpcURL))
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	response := helpers.JSONResponse{
		Error:   false,
		Message: result,
	}
	_ = tool.WriteJSON(writer, http.StatusAccepted, response)
}

func (app *Config) logEventViaGRPC(writer http.ResponseWriter, request *http.Request) {
	var tool helpers.Tool
	var requestPayload RequestPayload

	if err := tool.ReadJSONBody(writer, request, &requestPayload); err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s", app.GrpcURL),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}
	defer conn.Close()

	client := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		_ = tool.WriteErrorJSON(writer, err)
		return
	}

	response := helpers.JSONResponse{
		Error:   false,
		Message: "logged via grpc",
	}
	err = tool.WriteJSON(writer, http.StatusAccepted, response)
	if err != nil {
		return
	}
}
