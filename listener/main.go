package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/bhusal-rj/listner/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct{}

func main() {
	//try to connect to rabbitmq
	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	//start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages")

	consumer, err := event.NewConsumer(conn)
	if err != nil {
		log.Panic("There has been an erro while consuming rabbitmq", consumer)
		return
	}

	//watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panic("Error while listening to RabbitMQ", err)
		return
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
