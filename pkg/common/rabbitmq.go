package common

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func RetrieveRabbitMQQueue() (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	fmt.Println(rabbitMQURL)
	var connection *amqp.Connection
	var channel *amqp.Channel
	var err error

	for retries := 0; retries < 5; retries++ {
		connection, err = amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ (attempt %d): %v", retries+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		channel, err = connection.Channel()
		if err != nil {
			log.Printf("Failed to open channel (attempt %d): %v", retries+1, err)
			connection.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		queue, err := channel.QueueDeclare(
			"task-queue", // Queue name
			false,        // Durable
			false,        // Auto delete
			false,        // Exclusive
			false,        // No wait
			nil,          // Arguments
		)
		if err != nil {
			log.Printf("Failed to declare queue (attempt %d): %v", retries+1, err)
			channel.Close()
			connection.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		return connection, channel, queue
	}

	log.Fatalf("Failed to connect to RabbitMQ after 5 attempts: %v", err)
	return nil, nil, amqp.Queue{}
}
