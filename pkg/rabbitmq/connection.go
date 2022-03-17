package rabbitmq

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

	return connectRabbitMQ, nil
}

func NewChanel(rabbitConnection *amqp.Connection) (*amqp.Channel, error) {
	channelRabbitMQ, err := rabbitConnection.Channel()
	if err != nil {
		return nil, err
	}

	return channelRabbitMQ, nil
}

func Close(connection *amqp.Connection, channel *amqp.Channel) {
	defer func(connection *amqp.Connection) {
		err := connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)

	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(channel)
}
