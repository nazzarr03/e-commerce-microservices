package config

import (
	"os"

	"github.com/streadway/amqp"
)

var (
	RabbitMQConn *amqp.Connection
)

func ConnectRabbitMQ() {
	var err error
	RabbitMQConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic(err)
	}
}
