package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (receiver *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.Background(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error inserting in mongo: ", err)
		return err
	}

	*resp = "Processed payload via RPC:" + payload.Name

	return nil
}
