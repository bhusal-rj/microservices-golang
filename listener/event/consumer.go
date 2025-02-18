package event

import (
	//advanced messaging queue protocol
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	//
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	//communication in the particular channel is completely separate from the communication on another channel
	//connection creation creates a channel and termination terminates the channel
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// listen to the rabbit mq
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()
	//to ensure that queue exist before receiving the messages
	q, err := declareRandomQueue(ch)

	if err != nil {
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			err = json.Unmarshal(d.Body, &payload)
			if err != nil {
				fmt.Println("Unmarshaling error", err)
			}
			handlePayload(payload)
			// var payload Payload
		}
	}()
	log.Print("Waiting for messages [*] in the queue", q.Name, <-forever)
	return nil
}

func handlePayload(payload Payload) {
	fmt.Println("Payload", payload)
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		//authenticate

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	}
}

func logEvent(l Payload) error {
	jsonData, _ := json.Marshal(l)

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil

}
