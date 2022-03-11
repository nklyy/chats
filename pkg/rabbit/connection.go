package rabbit

import (
	"github.com/streadway/amqp"
	"log"
)

func NewConnection(amqpUrl string) (*amqp.Connection, error) {
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, err
	}
	defer func(connectRabbitMQ *amqp.Connection) {
		err := connectRabbitMQ.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connectRabbitMQ)

	return connectRabbitMQ, nil
}

func NewChanel(rabbitConnection *amqp.Connection) (*amqp.Channel, error) {
	channelRabbitMQ, err := rabbitConnection.Channel()
	if err != nil {
		return nil, err
	}
	defer func(channelRabbitMQ *amqp.Channel) {
		err := channelRabbitMQ.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(channelRabbitMQ)

	return channelRabbitMQ, nil
}
