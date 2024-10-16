package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/bhusal-rj/logger/cmd/api/logs"
	"github.com/bhusal-rj/logger/cmd/data"
	"google.golang.org/grpc"
)

type LogServer struct {
	//backward compatibility
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	//write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)

	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	res := &logs.LogResponse{Result: "Logged"}
	return res, nil

}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		fmt.Println("There has been an error listening to the gRPC", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Printf("gRPC server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		fmt.Println("There has been an error listening to the gRPC", err)
	}

}
