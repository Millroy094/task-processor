package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/millroy094/task-processor/pkg/common"
	"github.com/millroy094/task-processor/pkg/task"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client
var taskCollection *mongo.Collection

func updateTaskStatusAndFinishedAt(taskID int, status string, finishedAt time.Time) {
	filter := map[string]interface{}{"id": taskID}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"status":     status,
			"finishedAt": finishedAt,
		},
	}
	_, err := taskCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating task status and finishedAt: %v", err)
	}
}

func stopRetrying(task task.Task) {
	updateTaskStatusAndFinishedAt(task.ID, "failed", time.Now())
	log.Printf("Task %d exceeded max retries and is marked as failed.", task.ID)
}

func performTask(task task.Task) error {
	switch task.Type {
	case "email":
		return sendEmail(task)
	case "health_check":
		return performHealthCheck(task)
	default:
		return nil
	}
}

func processor(id int, messages <-chan amqp.Delivery, shutdownChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case message, ok := <-messages:
			if !ok {
				log.Printf("Worker %d: Stopping\n", id)
				return
			}

			var task task.Task
			err := json.Unmarshal(message.Body, &task)

			if err != nil {
				log.Printf("Worker %d: Failed to deserialize task: %v\n", id, err)
				message.Nack(false, false)
				continue
			}

			retryCount, ok := message.Headers["retryCount"].(int)
			if !ok {
				retryCount = 0
			}

			maxRetries, err := strconv.Atoi(os.Getenv("MAX_RETRIES"))
			if err != nil {
				maxRetries = 3
			}

			if err := performTask(task); err != nil {
				log.Printf("Worker %d: Task failed: %v\n", id, err)

				if retryCount < maxRetries {
					retryCount++
					message.Headers["retryCount"] = retryCount
					message.Nack(true, false)
				} else {
					stopRetrying(task)
					message.Ack(false)
				}
			} else {
				updateTaskStatusAndFinishedAt(task.ID, "completed", time.Now())
				fmt.Printf("Worker %d processing task ID %d of type %s\n", id, task.ID, task.Type)
				message.Ack(false)
			}
		case <-shutdownChan:
			log.Printf("Worker %d: Received shutdown signal\n", id)
			return
		}
	}
}

func main() {
	_, err := common.PrepareEnvironment([]string{"RABBITMQ_URL", "MONGODB_URL", "MAX_RETRIES"})

	if err != nil {
		log.Fatalf("Environment preparation failed: %v", err)
	}

	mongoClient = common.InitializeMongoDb()
	taskCollection = mongoClient.Database("task_manager").Collection("tasks")

	connection, channel, queue := common.RetrieveRabbitMQQueue()

	defer connection.Close()
	defer channel.Close()

	if err := channel.Qos(10, 0, false); err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	messages, err := channel.Consume(queue.Name, "", false, false, false, false, nil)

	common.FailOnError(err, "Failed to register a consumer")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	shutdownChan := make(chan struct{})
	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go processor(i, messages, shutdownChan, &wg)
	}

	<-stopChan
	log.Println("Shutdown signal received, stopping workers...")

	close(shutdownChan)

	wg.Wait()
	log.Println("All workers stopped. Exiting...")
}
