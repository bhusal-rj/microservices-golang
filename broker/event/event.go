package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// Exchange is a routing mechanism within RabbitMQ that takes message from producers and directs them to one or more queues
// Producer never sends message to the queues instead publish message to an exchange
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      //type:- topic routes message based on the pattern matchin in the routing key
		true,         //durable
		false,        //autodeleted
		false,        //internal
		false,        //no-wait
		nil,          //arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name
		false, //durable
		false, //delete
		true,  //exclusive
		false, //no-wait
		nil,   //arguments
	)
}
