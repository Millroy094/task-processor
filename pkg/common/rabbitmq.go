package common

import (
	"os"
	"github.com/streadway/amqp"
)

func RetrieveRabbitMQQueue() (*amqp.Connection, *amqp.Channel, amqp.Queue) {

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	
	connection, err := amqp.Dial(rabbitMQURL)
	FailOnError(err, "Failed to connect to RabbitMQ")

	channel, err := connection.Channel()
	FailOnError(err, "Failed to open a channel")

	queue, err := channel.QueueDeclare(
		"task-queue",
		false,
		false,
		false,
		false,
		nil,
	)

	FailOnError(err, "Failed to declare a queue")

	return connection, channel, queue

}
