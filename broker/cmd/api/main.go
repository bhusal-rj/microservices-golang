package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "3004"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	//connect to RabbitMQ
	conn, err := connect()
	if err != nil {
		log.Panic("Error connecting to the RabbitMQ server", err)
		return

	}
	defer conn.Close()

	log.Print("Successfully connected to the RabbitMQ")
	app := Config{Rabbit: conn}

	log.Printf("Starting the broker service at port %s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.route(),
	}

	//start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic("There has been an error", err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection = &amqp.Connection{}
	//donnot connect unless the rabbitmq is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not ready", err)
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("Backoff time reached its limit")
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)

	}

	return connection, nil
}
