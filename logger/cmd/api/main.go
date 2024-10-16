package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bhusal-rj/logger/cmd/data"
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
	mongoURL = "mongodb://admin:password@mongo:27017/logs?authSource=admin"
	gRpcPort = "50001"
)

// mongo client to connect to the mongoDB
var client *mongo.Client

type Config struct {
	Models data.Models
}

func (app *Config) rpcListen(wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", rpcPort))
	if err != nil {
		fmt.Println("There has been an error", err)
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			fmt.Println("There has been an error ", err)
		}
		go rpc.ServeConn(rpcConn)
	}
}

func main() {
	//connect to mongo
	mongoClient := connectToMongo()
	if mongoClient == nil {
		log.Panicf("Error connecting the database the program will exit")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	log.Print("Mongodb connection successful")
	app := Config{
		Models: data.New(mongoClient),
	}

	//Register the RPC server
	err := rpc.Register(new(RPCServer))

	if err != nil {
		fmt.Println("Error connecting to the RPC client", err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	//server for rpc
	go app.rpcListen(&wg)

	//server the logger REST server
	go app.serve(&wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Recevied interrupt signal")
	wg.Wait()
	log.Println("Servers shutdown")
}

func (app *Config) serve(wg *sync.WaitGroup) {
	defer wg.Done()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.route(),
	}
	log.Print("Logger service is up and running")
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

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
