package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const webPort = "80"
// const rpcPort = "5001"
// const mongoURL = "mongodb://mongo:27017"
// const gRpcPort = "50001"

const (
	webPort = "80"
	rpcPort = "5001"
	//connection_string in the form mongodb://username:password@host:port
	mongoURL = "mongodb://admin:password@mongo:27017"
	gRpcPort = "50001"
)

// mongo client to connect to the mongoDB
var client *mongo.Client

type Config struct {
}

func main() {
	//connect to mongo
	mongoClient := connectToMongo()
	if mongoClient == nil {
		log.Panicf("Error connecting the database the program will exit")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	log.Print("Mongodb connection successful")
}

func connectToMongo() *mongo.Client {
	opts := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		log.Panicf("There has been an error connecting to the database %v", err)
		return nil
	}

	return client

}
