package main

import (
	"os"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	"github.com/millroy094/task-processor/pkg/common"
	"github.com/millroy094/task-processor/pkg/task"
	_ "github.com/millroy094/task-processor/cmd/producer/docs"
)

// @title Task Processor API
// @version 1.0
// @description API for creating tasks and sending them to RabbitMQ.
// @host localhost:8080
// @BasePath /
func sendTask(channel *amqp.Channel, queue string, task task.Task) {
	body, err := json.Marshal(task)
	common.FailOnError(err, "Failed to serialize task")

	err = channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	common.FailOnError(err, "Failed to publish task")

	log.Printf("Task ID %d sent to queue %s\n", task.ID, queue)
}

// @Summary Create a new task
// @Description Create a new task and send it to the RabbitMQ queue.
// @Accept  json
// @Produce  json
// @Param task body task.Task true "Task" // Task definition for Swagger
// @Success 201 {object} map[string]interface{} "Task created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /tasks [post]
func createTaskHandler(channel *amqp.Channel, queue string) func(*gin.Context) {
	return func(c *gin.Context) { 
		var task task.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		sendTask(channel, queue, task)
		c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Task ID %d queued successfully", task.ID)})
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	apiPort := os.Getenv("API_PORT")

	if rabbitMQURL == "" || apiPort == "" {
		log.Fatal("Missing required environment variables")
	}
	
	connection, channel, queue := common.RetrieveRabbitMQQueue()

	defer connection.Close()
	defer channel.Close()

	r := gin.Default()


	r.POST("/tasks", createTaskHandler(channel, queue.Name))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + apiPort)
}
