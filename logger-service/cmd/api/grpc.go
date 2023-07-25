package main

import (
	"context"
	"time"

	"log-service/data"
	"log-service/logs"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, request *logs.LogRequest) (*logs.LogResponse, error) {
	input := request.GetLogEntry()

	logEntry := data.LogEntry{
		Name:      input.Name,
		Data:      input.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := l.Models.LogEntry.InsertRecord(logEntry); err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	res := &logs.LogResponse{Result: "logged"}
	return res, nil
}
