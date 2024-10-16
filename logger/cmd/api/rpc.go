package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bhusal-rj/logger/cmd/data"
)

// type of RPC server
type RPCServer struct {
}

// data accepted by the RPC server
type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	fmt.Println("Logging the result")

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Erroring writing to mongo", err)
		return err
	}

	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}
