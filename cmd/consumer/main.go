package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/millroy094/task-processor/pkg/common"
	"github.com/millroy094/task-processor/pkg/task"
	"github.com/streadway/amqp"
)

func performTask(task task.Task) error {
	// Simulate task handling logic
	// Return nil on success or an error on failure
	return nil
}

func processor(id int, messages <-chan amqp.Delivery) {
	for message := range messages {
		var task task.Task
		err := json.Unmarshal(message.Body, &task)

		if err != nil {
			log.Printf("Worker %d: Failed to deserialize task: %v\n", id, err)
			message.Nack(false, false)
			continue
		}

		if err := performTask(task); err != nil {
			log.Printf("Worker %d: Task failed: %v\n", id, err)
			message.Nack(false, true)
			continue
		}

		fmt.Printf("Worker %d processing task ID %d of type %s\n", id, task.ID, task.Type)

		message.Ack(false)
	}
}

func main() {

	_, err := common.PrepareEnvironment([]string{"RABBITMQ_URL"})

	if err != nil {
		log.Fatalf("Environment preparation failed: %v", err)
	}

	connection, channel, queue := common.RetrieveRabbitMQQueue()

	defer connection.Close()
	defer channel.Close()

	if err := channel.Qos(10, 0, false); err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	messages, err := channel.Consume(queue.Name, "", false, false, false, false, nil)

	common.FailOnError(err, "Failed to register a consumer")

	numWorkers := runtime.NumCPU()
	for i := 1; i <= numWorkers; i++ {
		go processor(i, messages)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
}
