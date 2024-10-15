package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//try to connect to rabbitmq
	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	//start listening for messages

	//watch the queue and cosume events
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection = &amqp.Connection{}
	//donnot connect unless the rabbitmq is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
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
