package producer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/millroy094/task-processor/pkg/common"
	"github.com/streadway/amqp"
)

func sendTask(channel *amqp.Channel, queue string, task common.Task) {
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

func main() {
	
	connection, channel, queue := common.RetrieveRabbitMQQueue()

	defer connection.Close()
	defer channel.Close()

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var task common.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		sendTask(channel, queue.Name, task)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Task ID %d queued successfully", task.ID)
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
