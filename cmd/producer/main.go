package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/millroy094/task-processor/pkg/common"
	"github.com/millroy094/task-processor/pkg/task"
	"github.com/streadway/amqp"
	v3 "github.com/swaggest/swgui/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client
var taskCollection *mongo.Collection

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
// @Param task body task.Task true "Task"
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

		task.CreatedAt = time.Now()
		task.Status = "pending"

		_, err := taskCollection.InsertOne(nil, task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task to MongoDB"})
			return
		}

		sendTask(channel, queue, task)
		c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Task ID %d queued successfully", task.ID)})
	}
}

// @Summary Get all tasks
// @Description Get tasks for all statuses.
// @Accept  json
// @Produce  json
// @Success 200 {array} task.Task
// @Failure 400 {object} map[string]interface{} "Failed to retrieve tasks from MongoDB"
// @Router /tasks [get]
func getAllTasksHandler() func(*gin.Context) {
	return func (c *gin.Context) {

		cursor, err := taskCollection.Find(context.TODO(), bson.D{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks from MongoDB"})
			return 
		}

		var results []task.Task

		if err = cursor.All(context.TODO(), &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks from MongoDB"})
		}

		c.JSON(http.status.Ok, results)

	}
}

func main() {

	envVariables, err := common.PrepareEnvironment([]string{"RABBITMQ_URL", "MONGODB_URL", "API_PORT"})

	if err != nil {
		log.Fatalf("Environment preparation failed: %v", err)
	}

	apiPort := envVariables["API_PORT"]

	mongoClient = common.InitializeMongoDb()
	taskCollection = mongoClient.Database("task_manager").Collection("tasks")

	connection, channel, queue := common.RetrieveRabbitMQQueue()

	defer connection.Close()
	defer channel.Close()

	r := gin.Default()

	r.POST("/tasks", createTaskHandler(channel, queue.Name))
	r.GET("/tasks", getAllTasksHandler())

	swaggerHandler := v3.NewHandler("Task Processor API", "/swagger.json", "/swagger/")
	r.GET("/swagger/*any", func(c *gin.Context) {
		swaggerHandler.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	r.Run(":" + apiPort)
}
